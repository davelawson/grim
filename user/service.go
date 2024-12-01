package user

import (
	"main/utils"
)

type Service struct {
	userRepo *Repo
}

func NewService(userRepo *Repo) *Service {
	instance := &Service{userRepo: userRepo}
	return instance
}

func (us *Service) CreateUser(email string, name string, password string) error {
	passwordHash, err := utils.Hash(password, email)
	if err != nil {
		return err
	}
	return us.userRepo.CreateUser(email, name, passwordHash)
}

func (us *Service) GetUserByEmail(email string) (*User, error) {
	user, err := us.userRepo.GetUserByEmail(email)
	return user, err
}
