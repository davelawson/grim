package match

import (
	"database/sql"
	"main/game"
	"main/util"
)

type ServiceFacade struct {
	service *Service
	db      *sql.DB
}

func NewServiceFacade(service *Service) *ServiceFacade {
	return &ServiceFacade{service: service, db: service.repo.db}
}

func (sf *ServiceFacade) Get(id, actorID string) (*game.View, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*game.View, error) { return sf.service.Get(tx, id, actorID) })()
}

func (sf *ServiceFacade) End(id, actorID string) error {
	return util.InTx(sf.db, func(tx *sql.Tx) error { return sf.service.End(tx, id, actorID) })()
}

func (sf *ServiceFacade) Command(id, actorID string, command CommandRequest) (*game.View, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*game.View, error) { return sf.service.Command(tx, id, actorID, command) })()
}

func (sf *ServiceFacade) Catalogue(id, actorID string) (*game.Catalogue, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*game.Catalogue, error) { return sf.service.Catalogue(tx, id, actorID) })()
}
