package auth

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"main/model"
	"main/util"
)

type userRepo interface {
	GetUserByEmail(tx *sql.Tx, email string) (*model.User, error)
	GetUserByToken(tx *sql.Tx, token string) (*model.User, error)
	UpdateUser(tx *sql.Tx, user *model.User) error
}

type Service struct {
	userRepo userRepo
}

func NewService(userRepo userRepo) *Service {
	return &Service{userRepo}
}

// Always use lower-case for emails
func (as *Service) Login(tx *sql.Tx, email string, password string) (*string, error) {
	user, err := as.userRepo.GetUserByEmail(tx, email)
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
	tokenAsString := base64.StdEncoding.EncodeToString(token)
	user.Token = tokenAsString
	as.userRepo.UpdateUser(tx, user)

	return &tokenAsString, nil
}

func (as *Service) VerifyBearerToken(tx *sql.Tx, token string) (*model.User, error) {
	user, err := as.userRepo.GetUserByToken(tx, token)
	if user == nil {
		fmt.Println("Unable to verify auth token ", token)
	} else {
		fmt.Println("Verified auth token ", token)
	}
	return user, err
}
