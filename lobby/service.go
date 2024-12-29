package lobby

import "errors"

type userRepo interface {
	GetUserByUUID(uuid string)
}

type Service struct {
	repo     *LobbyRepo
	userRepo userRepo
}

func (ls *Service) AddUserToLobby(lobbyId string, userId string) error {
	return ls.repo.AddUserToLobby(lobbyId, userId)
}

func NewService(repo *LobbyRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (ls *Service) CreateLobby(name string, userId string) (string, error) {
	lobbyId, err := ls.repo.CreateLobby(name, userId)
	if err != nil {
		return "", err
	}
	err = ls.repo.AddMemberToLobby(lobbyId, userId)
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
