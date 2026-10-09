package match

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/util"
	"strings"

	"github.com/google/uuid"
)

type adminRepo interface {
	IsAdmin(*sql.Tx, string) (bool, error)
}

type Service struct {
	repo  *Repo
	users adminRepo
}

func NewService(repo *Repo, users adminRepo) *Service { return &Service{repo: repo, users: users} }

// CreateSetup participates in the lobby's transaction; it never commits alone.
func (s *Service) CreateSetup(tx *sql.Tx, name, lobbyID string, participants []string, actorID string) (*game.View, error) {
	state, err := game.NewSetup(uuid.NewString(), name, lobbyID, participants)
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

func decode(stored *storedMatch) (*game.State, error) {
	state, err := game.Decode(stored.snapshot)
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
	state, err := decode(stored)
	if err != nil {
		return nil, err
	}
	return state.View(actorID), nil
}

func (s *Service) Command(tx *sql.Tx, id, actorID string, command api.CommandRequest) (*game.View, error) {
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
	state, err := decode(stored)
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

func (s *Service) Catalogue(tx *sql.Tx, id, actorID string) (*game.Catalogue, error) {
	// Get applies membership, deletion, metadata and version checks.
	if _, err := s.Get(tx, id, actorID); err != nil {
		return nil, err
	}
	return game.OpeningCatalogue(), nil
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
