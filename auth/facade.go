package auth

import (
	"database/sql"
	"main/model"
	"main/util"
)

type ServiceFacade struct {
	service *Service
	db      *sql.DB
}

func NewServiceFacade(authService *Service, db *sql.DB) *ServiceFacade {
	return &ServiceFacade{
		service: authService,
		db:      db,
	}
}

func (sf *ServiceFacade) Login(email string, password string) (*string, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*string, error) {
		return sf.service.Login(tx, email, password)
	})()
}

func (sf *ServiceFacade) VerifyBearerToken(token string) (*model.User, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*model.User, error) {
		return sf.service.VerifyBearerToken(tx, token)
	})()
}
