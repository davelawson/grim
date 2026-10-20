package match

import (
	"database/sql"
	"log/slog"

	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/util"
)

type ServiceFacade struct {
	service *Service
	db      *sql.DB
}

func NewServiceFacade(service *Service) *ServiceFacade {
	return &ServiceFacade{service: service, db: service.repo.db}
}

func (sf *ServiceFacade) Get(id, actorID string) (*game.View, error) {
	return util.InTxWithResult(sf.db, func(tx *sql.Tx) (*game.View, error) { return sf.service.Get(tx, id, actorID) })()
}

func (sf *ServiceFacade) End(id, actorID string) error {
	return util.InTx(sf.db, func(tx *sql.Tx) error { return sf.service.End(tx, id, actorID) })()
}

func (sf *ServiceFacade) GetDraft(id, actorID string) (*api.DraftResponse, error) {
	return util.InTxWithResult(sf.db, func(tx *sql.Tx) (*api.DraftResponse, error) { return sf.service.GetDraft(tx, id, actorID) })()
}

func (sf *ServiceFacade) SubmitDraft(id, actorID string, command api.DraftPickRequest) (*api.DraftResponse, error) {
	return util.InTxWithResult(sf.db, func(tx *sql.Tx) (*api.DraftResponse, error) { return sf.service.SubmitDraft(tx, id, actorID, command) })()
}

func (sf *ServiceFacade) Catalogue(id, actorID string) (*game.Catalogue, error) {
	return util.InTxWithResult(sf.db, func(tx *sql.Tx) (*game.Catalogue, error) { return sf.service.Catalogue(tx, id, actorID) })()
}

func (sf *ServiceFacade) submitGameplay(id, actorID string, command game.Command, operation gameplayOperation) (*game.View, error) {
	diagnostics := &commandDiagnostics{}
	view, err := util.InTxWithResult(sf.db, func(tx *sql.Tx) (*game.View, error) {
		return sf.service.submitGameplay(tx, id, actorID, command, operation, diagnostics)
	})()
	if err == nil {
		slog.Debug("match gameplay input", "match", id, "actor", actorID, "operation", operation, "revision", view.Revision, "rules", view.RulesVersion, "catalogue", view.CatalogueVersion, "command", command, "replayed", diagnostics.Replayed, "random", diagnostics.Random, "choice", view.Choice)
	}
	return view, err
}
