package user

import (
	"database/sql"
	"main/model"
	"main/util"
)

type ServiceFacade struct {
	service *Service
	db      *sql.DB
}

func NewServiceFacade(service *Service) *ServiceFacade {
	return &ServiceFacade{
		service: service,
		db:      service.userRepo.db,
	}
}

func (sf *ServiceFacade) CreateUser(email string, name string, password string) error {
	return util.InTx(sf.db, func(tx *sql.Tx) error {
		return sf.service.CreateUser(tx, email, name, password)
	})()
}

func (sf *ServiceFacade) GetUserByEmail(email string) (*model.User, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*model.User, error) {
		return sf.service.GetUserByEmail(tx, email)
	})()
}

func (sf *ServiceFacade) GetUserByToken(token string) (*model.User, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*model.User, error) {
		return sf.service.GetUserByToken(tx, token)
	})()
}
