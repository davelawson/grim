package lobby

import (
	"database/sql"
	"main/util"
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
	return util.InTx(sf.db, func(tx *sql.Tx) error {
		return sf.service.AddUserToLobby(tx, lobbyId, userId, requestorId)
	})()
}

func (sf *ServiceFacade) RemoveUserFromLobby(lobbyId string, userId string, requestorId string) error {
	return util.InTx(sf.db, func(tx *sql.Tx) error {
		return sf.service.RemoveUserFromLobby(tx, lobbyId, userId, requestorId)
	})()
}

func (sf *ServiceFacade) UpdateLobby(lobbyId string, name string, ownerId string) error {
	return util.InTx(sf.db, func(tx *sql.Tx) error {
		return sf.service.UpdateLobby(tx, lobbyId, name, ownerId)
	})()
}

func (sf *ServiceFacade) DeleteLobby(lobbyId string, userId string) error {
	return util.InTx(sf.db, func(tx *sql.Tx) error {
		return sf.service.DeleteLobby(tx, lobbyId, userId)
	})()
}

func (sf *ServiceFacade) CreateLobby(name string, userId string) (*string, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*string, error) {
		return sf.service.CreateLobby(tx, name, userId)
	})()
}

func (sf *ServiceFacade) GetLobby(id string) (*Lobby, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*Lobby, error) {
		return sf.service.GetLobby(tx, id)
	})()
}
