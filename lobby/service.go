package lobby

import (
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
	NotOwner:           errors.New("only the owner can update a lobby"),
	UserAlreadyInLobby: errors.New("user already in lobby"),
}

func (ls *Service) AddUserToLobby(lobbyId string, userId string, requestorId string) error {
	lobby, err := ls.GetLobby(lobbyId)
	if err != nil {
		return err
	} else if lobby == nil {
		return AddUserToLobbyErrors.LobbyNotFound
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

	return ls.repo.AddUserToLobby(lobbyId, userId)
}

type UpdateLobbyErrorsType struct {
	NotFound error
}

var UpdateLobbyErrors = UpdateLobbyErrorsType{
	NotFound: errors.New("user does now own a lobby with the specified id"),
}

func (ls *Service) UpdateLobby(lobbyId string, name string, ownerId string) error {
	rowsAffected, err := ls.repo.UpdateLobby(lobbyId, name, ownerId)
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return UpdateLobbyErrors.NotFound
	}
	return nil
}

func (ls *Service) CreateLobby(name string, userId string) (string, error) {
	lobbyId, err := ls.repo.CreateLobby(name, userId)
	if err != nil {
		return "", err
	}
	err = ls.repo.AddUserToLobby(lobbyId, userId)
	if err != nil {
		return "", err
	}
	return lobbyId, err
}

type DeleteLobbyErrorsType struct {
	NotFound error
}

var DeleteLobbyErrors = DeleteLobbyErrorsType{
	NotFound: errors.New("user does not own a lobby with the specified id"),
}

func (ls *Service) DeleteLobby(lobbyId string, userId string) error {
	rows, err := ls.repo.DeleteLobby(lobbyId, userId)
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

func (ls *Service) GetLobby(id string) (*Lobby, error) {
	lobby, err := ls.repo.GetLobby(id)
	if err != nil {
		return nil, err
	} else if lobby == nil {
		return nil, GetLobbyErrors.NotFound
	}

	members, err := ls.repo.GetLobbyMembers(id)
	if err != nil {
		return nil, err
	}
	lobby.Members = members
	return lobby, err
}
