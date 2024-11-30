package service

import (
	"main/model"
	"main/repo"
)

type UserService struct {
	userRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	instance := &UserService{userRepo: userRepo}
	return instance
}

func (us *UserService) CreateUser(email string, name string, password string) error {
	passwordHash, err := Hash(password, email)
	if err != nil {
		return err
	}
	return us.userRepo.CreateUser(email, name, passwordHash)
}

func (us *UserService) GetUserByEmail(email string) (*model.User, error) {
	user, err := us.userRepo.GetUserByEmail(email)
	return user, err
}
