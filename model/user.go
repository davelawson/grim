package model

type User struct {
	Id           string
	Name         string
	Email        string
	PasswordHash []byte
	Token        string
}
