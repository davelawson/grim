package lobby

import (
	"database/sql"
	"errors"
	"main/game"
	"main/match"
	"main/model"
	"main/util"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type userRepo interface {
	GetUserById(tx *sql.Tx, id string) (*model.User, error)
}

type Service struct {
	repo     *LobbyRepo
	userRepo userRepo
	matches  *match.Service
}

func NewService(repo *LobbyRepo, userRepo userRepo, matches *match.Service) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
		matches:  matches,
	}
}

type AddUserToLobbyErrorsType struct {
	LobbyNotFound      error
	UserNotFound       error
	NotOwner           error
	UserAlreadyInLobby error
}

var AddUserToLobbyErrors = AddUserToLobbyErrorsType{
	LobbyNotFound:      errors.New("lobby not found"),
	UserNotFound:       errors.New("user not found"),
	NotOwner:           errors.New("only the owner can add a user to a lobby"),
	UserAlreadyInLobby: errors.New("user already in lobby"),
}

func (ls *Service) AddUserToLobby(tx *sql.Tx, lobbyId string, userId string, requestorId string) error {
	lobby, err := ls.GetLobby(tx, lobbyId)
	if err == GetLobbyErrors.NotFound {
		return AddUserToLobbyErrors.LobbyNotFound
	} else if err != nil {
		return err
	} else if lobby.Owner != requestorId {
		return AddUserToLobbyErrors.NotOwner
	} else if lobby.Status != "open" {
		return util.ErrLobbyClosed
	} else if slices.Contains(lobby.Members, userId) {
		return AddUserToLobbyErrors.UserAlreadyInLobby
	}

	user, err := ls.userRepo.GetUserById(tx, userId)
	if err != nil {
		return err
	} else if user == nil {
		return AddUserToLobbyErrors.UserNotFound
	}

	if err := ls.repo.AddUserToLobby(tx, lobbyId, userId); err != nil {
		return err
	}
	return ls.repo.ClearReadiness(tx, lobbyId)
}

type RemoveUserFromLobbyErrorsType struct {
	LobbyNotFound  error
	NotOwner       error
	UserNotInLobby error
}

var RemoveUserFromLobbyErrors = RemoveUserFromLobbyErrorsType{
	LobbyNotFound:  errors.New("lobby not found"),
	NotOwner:       errors.New("only the owner can remove a user from a lobby"),
	UserNotInLobby: errors.New("user not in lobby"),
}

func (ls *Service) RemoveUserFromLobby(tx *sql.Tx, lobbyId string, userId string, requestorId string) error {
	lobby, err := ls.GetLobby(tx, lobbyId)
	if err == GetLobbyErrors.NotFound {
		return RemoveUserFromLobbyErrors.LobbyNotFound
	} else if err != nil {
		return err
	} else if lobby.Owner != requestorId && userId != requestorId {
		return RemoveUserFromLobbyErrors.NotOwner
	} else if !slices.Contains(lobby.Members, requestorId) {
		return util.ErrForbidden
	} else if lobby.Status != "open" {
		return util.ErrLobbyClosed
	} else if !slices.Contains(lobby.Members, userId) {
		return RemoveUserFromLobbyErrors.UserNotInLobby
	} else if userId == lobby.Owner {
		return &util.APIError{Status: 409, Code: "owner_must_transfer", Message: "Transfer ownership or delete the open lobby before leaving"}
	}

	if err := ls.repo.RemoveMemberFromLobby(tx, lobbyId, userId); err != nil {
		return err
	}
	return ls.repo.ClearReadiness(tx, lobbyId)
}

type UpdateLobbyErrorsType struct {
	NotFound error
}

var UpdateLobbyErrors = UpdateLobbyErrorsType{
	NotFound: util.ErrNotFound,
}

func (ls *Service) UpdateLobby(tx *sql.Tx, lobbyId string, name string, ownerId string, requestorId string) error {
	lobby, err := ls.GetLobby(tx, lobbyId)
	if err != nil {
		return err
	}
	if lobby.Owner != requestorId || !slices.Contains(lobby.Members, requestorId) {
		return util.ErrForbidden
	}
	if lobby.Status != "open" {
		return util.ErrLobbyClosed
	}
	if !slices.Contains(lobby.Members, ownerId) {
		return &util.APIError{Status: 400, Code: "invalid_owner", Message: "The new owner must be a lobby member"}
	}
	rowsAffected, err := ls.repo.UpdateLobby(tx, lobbyId, name, ownerId)
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return UpdateLobbyErrors.NotFound
	}
	return nil
}

func (ls *Service) CreateLobby(tx *sql.Tx, name string, userId string) (*string, error) {
	lobbyId, err := ls.repo.CreateLobby(tx, name, userId)
	if err != nil {
		return nil, err
	}
	err = ls.repo.AddUserToLobby(tx, lobbyId, userId)
	if err != nil {
		return nil, err
	}
	return &lobbyId, err
}

type DeleteLobbyErrorsType struct {
	NotFound error
}

var DeleteLobbyErrors = DeleteLobbyErrorsType{
	NotFound: util.ErrNotFound,
}

func (ls *Service) DeleteLobby(tx *sql.Tx, lobbyId string, userId string) error {
	lobby, err := ls.GetLobby(tx, lobbyId)
	if err != nil {
		return err
	}
	if lobby.Owner != userId {
		return util.ErrForbidden
	}
	if lobby.Status != "open" {
		return util.ErrLobbyClosed
	}
	rows, err := ls.repo.DeleteLobby(tx, lobbyId, userId)
	if err != nil {
		return err
	}
	if rows == 0 {
		return DeleteLobbyErrors.NotFound
	}
	return nil
}

type GetLobbyErrorsType struct {
	NotFound error
}

var GetLobbyErrors = GetLobbyErrorsType{
	NotFound: util.ErrNotFound,
}

func (ls *Service) GetLobby(tx *sql.Tx, id string) (*Lobby, error) {
	lobby, err := ls.repo.GetLobby(tx, id)
	if err != nil {
		return nil, err
	} else if lobby == nil {
		return nil, GetLobbyErrors.NotFound
	}

	members, err := ls.repo.GetLobbyMembers(tx, id)
	if err != nil {
		return nil, err
	}
	lobby.Members = members
	lobby.Readiness, err = ls.repo.GetReadiness(tx, id)
	return lobby, err
}

func (ls *Service) GetLobbyForUser(tx *sql.Tx, id, requestorId string) (*Lobby, error) {
	lobby, err := ls.GetLobby(tx, id)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(lobby.Members, requestorId) {
		return nil, util.ErrForbidden
	}
	return lobby, nil
}

func (ls *Service) SetReady(tx *sql.Tx, id, requestorId string, ready bool) (*Lobby, error) {
	lobby, err := ls.GetLobbyForUser(tx, id, requestorId)
	if err != nil {
		return nil, err
	}
	if lobby.Status != "open" {
		return nil, util.ErrLobbyClosed
	}
	if err := ls.repo.SetReady(tx, id, requestorId, ready); err != nil {
		return nil, err
	}
	lobby.Readiness[requestorId] = ready
	return lobby, nil
}

func (ls *Service) Launch(tx *sql.Tx, id, actorID, requestID string) (*game.View, error) {
	// Validate for non-HTTP callers as well as the controller.
	parsed, err := uuid.Parse(requestID)
	if err != nil || !strings.EqualFold(requestID, parsed.String()) {
		return nil, util.ErrInvalidRequest
	}
	requestID = parsed.String()
	lobby, err := ls.GetLobbyForUser(tx, id, actorID)
	if err != nil {
		return nil, err
	}
	if lobby.Owner != actorID {
		return nil, util.ErrForbidden
	}
	if view, err := ls.matches.LaunchReceipt(tx, id, actorID, requestID); err != nil || view != nil {
		return view, err
	}
	if lobby.Status != "open" {
		return nil, util.ErrLobbyClosed
	}
	if len(lobby.Members) < 2 || len(lobby.Members) > 4 {
		return nil, &util.APIError{Status: 409, Code: "invalid_roster_size", Message: "Launch requires two to four members"}
	}
	for _, member := range lobby.Members {
		if !lobby.Readiness[member] {
			return nil, &util.APIError{Status: 409, Code: "lobby_not_ready", Message: "Every member must be ready before launch"}
		}
	}
	view, err := ls.matches.CreateSetup(tx, lobby.Name, id, lobby.Members, actorID)
	if err != nil {
		return nil, err
	}
	if err := ls.repo.Close(tx, id, view.ID); err != nil {
		return nil, err
	}
	if err := ls.matches.SaveLaunchReceipt(tx, id, actorID, requestID, view); err != nil {
		return nil, err
	}
	return view, nil
}
