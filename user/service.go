package user

import (
	"main/model"
	"main/util"
)

type userRepo interface {
	CreateUser(email string, name string, passwordHash []byte) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByToken(token string) (*model.User, error)
}

type Service struct {
	userRepo userRepo
}

func NewService(userRepo userRepo) *Service {
	return &Service{userRepo: userRepo}
}

func (us *Service) CreateUser(email string, name string, password string) error {
	passwordHash, err := util.Hash(password, email)
	if err != nil {
		return err
	}
	return us.userRepo.CreateUser(email, name, passwordHash)
}

func (us *Service) GetUserByEmail(email string) (*model.User, error) {
	user, err := us.userRepo.GetUserByEmail(email)
	return user, err
}

func (us *Service) GetUserByToken(token string) (*model.User, error) {
	user, err := us.userRepo.GetUserByToken(token)
	return user, err
}
