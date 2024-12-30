package lobby

import (
	"database/sql"
	"errors"
	"main/model"
	"slices"
)

type userRepo interface {
	GetUserById(id string) (*model.User, error)
}

type Service struct {
	repo     *LobbyRepo
	userRepo userRepo
}

func NewService(repo *LobbyRepo, userRepo userRepo) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
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
	} else if slices.Contains(lobby.Members, userId) {
		return AddUserToLobbyErrors.UserAlreadyInLobby
	}

	user, err := ls.userRepo.GetUserById(userId)
	if err != nil {
		return err
	} else if user == nil {
		return AddUserToLobbyErrors.UserNotFound
	}

	return ls.repo.AddUserToLobby(tx, lobbyId, userId)
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
	} else if lobby.Owner != requestorId {
		return RemoveUserFromLobbyErrors.NotOwner
	} else if !slices.Contains(lobby.Members, userId) {
		return RemoveUserFromLobbyErrors.UserNotInLobby
	}

	return ls.repo.RemoveMemberFromLobby(tx, lobbyId, userId)
}

type UpdateLobbyErrorsType struct {
	NotFound error
}

var UpdateLobbyErrors = UpdateLobbyErrorsType{
	NotFound: errors.New("user does now own a lobby with the specified id"),
}

func (ls *Service) UpdateLobby(tx *sql.Tx, lobbyId string, name string, ownerId string) error {
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
	NotFound: errors.New("user does not own a lobby with the specified id"),
}

func (ls *Service) DeleteLobby(tx *sql.Tx, lobbyId string, userId string) error {
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
	NotFound: errors.New("no lobby found with the specified id"),
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
	return lobby, err
}
