package service

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"main/model"
)

type userRepo interface {
	GetUserByEmail(email string) (*model.User, error)
}

type AuthService struct {
	userRepo userRepo
}

func NewAuthService(userRepo userRepo) *AuthService {
	return &AuthService{userRepo}
}

// Always use lower-case for emails
func (as *AuthService) Login(email string, password string) ([]byte, error) {
	user, err := as.userRepo.GetUserByEmail(email)
	if err != nil {
		// TODO: Bubble up the error -- should probably result in an InternalServerError
		fmt.Println("Error getting user by email: {} -> {}", err, email)
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	hash, hashErr := Hash(password, email)
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

func (as *AuthService) VerifyBearerToken(token string) (*model.User, error) {
	return nil, nil
}
