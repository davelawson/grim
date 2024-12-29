package lobby

import (
	"database/sql"
)

type ServiceFacade struct {
	service *Service
	db      *sql.DB
}

func NewServiceFacade(lobbyService *Service) *ServiceFacade {
	return &ServiceFacade{
		service: lobbyService,
		db:      lobbyService.repo.db,
	}
}

func (sf *ServiceFacade) AddUserToLobby(lobbyId string, userId string, requestorId string) error {
	return sf.service.AddUserToLobby(lobbyId, userId, requestorId)
}

func (sf *ServiceFacade) RemoveUserFromLobby(lobbyId string, userId string, requestorId string) error {
	return sf.service.RemoveUserFromLobby(lobbyId, userId, requestorId)
}

func (sf *ServiceFacade) UpdateLobby(lobbyId string, name string, ownerId string) error {
	return sf.service.UpdateLobby(lobbyId, name, ownerId)
}

func (sf *ServiceFacade) DeleteLobby(lobbyId string, userId string) error {
	return sf.service.DeleteLobby(lobbyId, userId)
}

func (sf *ServiceFacade) CreateLobby(name string, userId string) (string, error) {
	return sf.service.CreateLobby(name, userId)
}

func (sf *ServiceFacade) GetLobby(id string) (*Lobby, error) {
	return sf.service.GetLobby(id)
}
