package auth

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"main/model"
	"main/util"
)

type userRepo interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserByToken(token string) (*model.User, error)
}

type Service struct {
	userRepo userRepo
}

func NewService(userRepo userRepo) *Service {
	return &Service{userRepo}
}

// Always use lower-case for emails
func (as *Service) Login(email string, password string) ([]byte, error) {
	user, err := as.userRepo.GetUserByEmail(email)
	if err != nil {
		// TODO: Bubble up the error -- should probably result in an InternalServerError
		fmt.Println("Error getting user by email: ", err, " -> ", email)
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	hash, hashErr := util.Hash(password, email)
	if hashErr != nil {
		return nil, hashErr
	}

	if !bytes.Equal(hash, user.PasswordHash) {
		fmt.Println("Hash doesn't match!")
		return nil, nil
	}

	// TODO: Generate a new bearer authentication token, and store that in the record
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (as *Service) VerifyBearerToken(token string) (*model.User, error) {
	user, err := as.userRepo.GetUserByToken(token)
	if user == nil {
		fmt.Println("Unable to verify auth token ", token)
	} else {
		fmt.Println("Verified auth token ", token)
	}
	return user, err
}
