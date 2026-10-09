package lobby

import (
	"database/sql"
	"main/api"
	"main/game"
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

func (sf *ServiceFacade) UpdateLobby(lobbyId string, name string, ownerId string, requestorId string) error {
	return util.InTx(sf.db, func(tx *sql.Tx) error {
		return sf.service.UpdateLobby(tx, lobbyId, name, ownerId, requestorId)
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

func (sf *ServiceFacade) GetLobby(id string, requestorId string) (*api.Lobby, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*api.Lobby, error) {
		return sf.service.GetLobbyForUser(tx, id, requestorId)
	})()
}

func (sf *ServiceFacade) SetReady(id, requestorId string, ready bool) (*api.Lobby, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*api.Lobby, error) { return sf.service.SetReady(tx, id, requestorId, ready) })()
}

func (sf *ServiceFacade) Launch(id, actorID, requestID string) (*game.View, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*game.View, error) { return sf.service.Launch(tx, id, actorID, requestID) })()
}
