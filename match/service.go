package match

import (
	"database/sql"
	"errors"
	"main/game"
	"main/util"

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
func (s *Service) CreateSetup(tx *sql.Tx, name, lobbyID string, participants []string) (*game.View, error) {
	state, err := game.NewSetup(uuid.NewString(), name, lobbyID, participants)
	if err != nil {
		return nil, err
	}
	if err := s.repo.create(tx, state); err != nil {
		return nil, err
	}
	return state.View(), nil
}

func (s *Service) LaunchReceipt(tx *sql.Tx, lobbyID, actorID, requestID string) (*game.View, error) {
	return s.repo.launchReceipt(tx, lobbyID, actorID, requestID)
}

func (s *Service) SaveLaunchReceipt(tx *sql.Tx, lobbyID, actorID, requestID string, view *game.View) error {
	return s.repo.saveLaunchReceipt(tx, lobbyID, actorID, requestID, view)
}

func decode(stored *storedMatch) (*game.State, error) {
	state, err := game.Decode(stored.snapshot)
	if errors.Is(err, game.ErrUnsupportedSchema) {
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
	return state.View(), nil
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
