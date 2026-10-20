package match

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/catalogue"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/util"
	"strings"

	"github.com/google/uuid"
)

type adminRepo interface {
	IsAdmin(*sql.Tx, string) (bool, error)
}

type Service struct {
	repo             *Repo
	users            adminRepo
	catalogueVersion string
}

func NewService(repo *Repo, users adminRepo, catalogueVersion string) *Service {
	return &Service{repo: repo, users: users, catalogueVersion: catalogueVersion}
}

// CreateSetup participates in the lobby's transaction; it never commits alone.
func (s *Service) CreateSetup(tx *sql.Tx, name, lobbyID string, participants []string, actorID string) (*game.View, error) {
	definitions, err := s.definitions(tx, s.catalogueVersion)
	if err != nil {
		return nil, err
	}
	state, err := game.NewSetup(uuid.NewString(), name, lobbyID, participants, definitions)
	if err != nil {
		return nil, err
	}
	if err := s.repo.create(tx, state); err != nil {
		return nil, err
	}
	return state.View(actorID), nil
}

func (s *Service) LaunchReceipt(tx *sql.Tx, lobbyID, actorID, requestID string) (*game.View, error) {
	return s.repo.launchReceipt(tx, lobbyID, actorID, requestID)
}

func (s *Service) SaveLaunchReceipt(tx *sql.Tx, lobbyID, actorID, requestID string, view *game.View) error {
	return s.repo.saveLaunchReceipt(tx, lobbyID, actorID, requestID, view)
}

func (s *Service) decode(tx *sql.Tx, stored *storedMatch) (*game.State, error) {
	var header struct {
		SchemaVersion    int    `json:"schemaversion"`
		RulesVersion     string `json:"rulesversion"`
		CatalogueVersion string `json:"catalogueversion"`
	}
	if err := json.Unmarshal(stored.snapshot, &header); err != nil {
		return nil, err
	}
	if header.SchemaVersion != game.SchemaVersion || header.RulesVersion != game.OpeningRulesVersion {
		return nil, util.ErrUnsupportedState
	}
	definitions, err := s.definitions(tx, header.CatalogueVersion)
	if err != nil {
		return nil, err
	}
	state, err := game.Decode(stored.snapshot, definitions)
	if errors.Is(err, game.ErrUnsupportedSchema) || errors.Is(err, game.ErrUnsupportedVersion) {
		return nil, util.ErrUnsupportedState
	}
	if err != nil {
		return nil, err
	}
	if state.ID != stored.id || state.Name != stored.name || state.LobbyID != stored.lobbyID ||
		state.Status != stored.status || state.Revision != stored.revision ||
		(state.Outcome != nil) != stored.outcome.Valid ||
		(state.Outcome != nil && *state.Outcome != stored.outcome.String) {
		return nil, errors.New("match metadata does not agree with saved state")
	}
	return state, nil
}

func (s *Service) Get(tx *sql.Tx, id, actorID string) (*game.View, error) {
	stored, err := s.repo.load(tx, id)
	if err != nil {
		return nil, err
	}
	member, err := s.repo.isParticipant(tx, id, actorID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, util.ErrForbidden
	}
	state, err := s.decode(tx, stored)
	if err != nil {
		return nil, err
	}
	return state.View(actorID), nil
}

func (s *Service) GetDraft(tx *sql.Tx, id, actorID string) (*api.DraftResponse, error) {
	view, err := s.Get(tx, id, actorID)
	if err != nil {
		return nil, err
	}
	return s.draftResponse(tx, view)
}

func (s *Service) SubmitDraft(tx *sql.Tx, id, actorID string, command api.DraftPickRequest) (*api.DraftResponse, error) {
	view, err := s.command(tx, id, actorID, command)
	if err != nil {
		return nil, err
	}
	return s.draftResponse(tx, view)
}

func (s *Service) draftResponse(tx *sql.Tx, view *game.View) (*api.DraftResponse, error) {
	definitions, err := s.definitions(tx, view.CatalogueVersion)
	if err != nil {
		return nil, err
	}
	response, err := game.DraftView(view, definitions)
	if errors.Is(err, game.ErrUnsupportedVersion) {
		return nil, util.ErrUnsupportedState
	}
	return response, err
}

// Keep participant views in receipts to support existing saved responses.
func (s *Service) command(tx *sql.Tx, id, actorID string, command api.DraftPickRequest) (*game.View, error) {
	parsed, err := uuid.Parse(command.RequestID)
	if err != nil || !strings.EqualFold(command.RequestID, parsed.String()) || command.ExpectedRevision == nil ||
		*command.ExpectedRevision < 0 || command.CardID == "" {
		return nil, util.ErrInvalidRequest
	}
	switch command.Type {
	case "choose_chantry", "choose_victory_path", "choose_minion":
	default:
		return nil, util.ErrInvalidRequest
	}
	command.RequestID = parsed.String()
	stored, err := s.repo.load(tx, id)
	if err != nil {
		return nil, err
	}
	member, err := s.repo.isParticipant(tx, id, actorID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, util.ErrForbidden
	}
	// Canonical typed input makes whitespace/key ordering irrelevant to retries.
	input, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}
	if view, err := s.repo.commandReceipt(tx, id, actorID, command.RequestID, string(input)); err != nil || view != nil {
		return view, err
	}
	state, err := s.decode(tx, stored)
	if err != nil {
		return nil, err
	}
	if state.Revision != *command.ExpectedRevision {
		return nil, util.ErrStaleRevision
	}
	next, err := state.Choose(actorID, command.Type, command.CardID)
	switch {
	case errors.Is(err, game.ErrNotParticipant):
		return nil, util.ErrForbidden
	case errors.Is(err, game.ErrDraftClosed):
		return nil, &util.APIError{Status: 409, Code: "draft_closed", Message: "Drafting has finished"}
	case errors.Is(err, game.ErrDraftSequence):
		return nil, &util.APIError{Status: 409, Code: "draft_sequence", Message: "This choice is not the current draft step"}
	case errors.Is(err, game.ErrInvalidChoice):
		return nil, &util.APIError{Status: 400, Code: "invalid_choice", Message: "This card is not offered for the current choice"}
	case errors.Is(err, game.ErrUnsupportedVersion):
		return nil, util.ErrUnsupportedState
	case err != nil:
		return nil, err
	}
	if err := s.repo.save(tx, next, state.Revision); err != nil {
		return nil, err
	}
	view := next.View(actorID)
	if err := s.repo.saveCommandReceipt(tx, id, actorID, command.RequestID, string(input), view); err != nil {
		return nil, err
	}
	return view, nil
}

// definitions loads only the requested version, never a fallback catalogue.
func (s *Service) definitions(tx *sql.Tx, version string) (*game.Definitions, error) {
	c, err := catalogue.Load(tx, version)
	if errors.Is(err, catalogue.ErrNotFound) {
		return nil, util.ErrUnsupportedState
	}
	if err != nil {
		return nil, err
	}
	definitions, err := game.NewDefinitions(*c)
	if errors.Is(err, game.ErrUnsupportedVersion) {
		return nil, util.ErrUnsupportedState
	}
	return definitions, err
}

// ValidateCatalogue checks startup configuration without installing or modifying data.
func (s *Service) ValidateCatalogue(tx *sql.Tx) error {
	_, err := s.definitions(tx, s.catalogueVersion)
	return err
}

func (s *Service) Catalogue(tx *sql.Tx, id, actorID string) (*game.Catalogue, error) {
	view, err := s.Get(tx, id, actorID)
	if err != nil {
		return nil, err
	}
	definitions, err := s.definitions(tx, view.CatalogueVersion)
	if err != nil {
		return nil, err
	}
	return definitions.Catalogue(), nil
}

func (s *Service) End(tx *sql.Tx, id, actorID string) error {
	admin, err := s.users.IsAdmin(tx, actorID)
	if err != nil {
		return err
	}
	if !admin {
		return util.ErrForbidden
	}
	return s.repo.softDelete(tx, id)
}
