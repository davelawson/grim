package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"main/model"
	"main/util"
)

type userRepo interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserByToken(token string) (*model.User, error)
	UpdateUser(user *model.User) error
}

type Service struct {
	userRepo userRepo
}

func NewService(userRepo userRepo) *Service {
	return &Service{userRepo}
}

// Always use lower-case for emails
func (as *Service) Login(email string, password string) (string, error) {
	user, err := as.userRepo.GetUserByEmail(email)
	if err != nil {
		// TODO: Bubble up the error -- should probably result in an InternalServerError
		fmt.Println("Error getting user by email: ", err, " -> ", email)
		return "", err
	}

	if user == nil {
		return "", nil
	}

	hash, hashErr := util.Hash(password, email)
	if hashErr != nil {
		return "", hashErr
	}

	if !bytes.Equal(hash, user.PasswordHash) {
		fmt.Println("Hash doesn't match!")
		return "", nil
	}

	// TODO: Generate a new bearer authentication token, and store that in the record
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		return "", err
	}
	tokenAsString := base64.StdEncoding.EncodeToString(token)
	user.Token = tokenAsString
	as.userRepo.UpdateUser(user)

	return tokenAsString, nil
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
