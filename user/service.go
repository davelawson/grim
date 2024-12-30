package user

import (
	"database/sql"
	"main/model"
	"main/util"
)

type Service struct {
	userRepo *UserRepo
}

func NewService(userRepo *UserRepo) *Service {
	return &Service{userRepo: userRepo}
}

func (us *Service) CreateUser(tx *sql.Tx, email string, name string, password string) error {
	passwordHash, err := util.Hash(password, email)
	if err != nil {
		return err
	}
	return us.userRepo.CreateUser(tx, email, name, passwordHash)
}

func (us *Service) GetUserByEmail(tx *sql.Tx, email string) (*model.User, error) {
	user, err := us.userRepo.GetUserByEmail(tx, email)
	return user, err
}

func (us *Service) GetUserByToken(tx *sql.Tx, token string) (*model.User, error) {
	user, err := us.userRepo.GetUserByToken(tx, token)
	return user, err
}
