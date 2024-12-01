package user

type User struct {
	Id           string
	Name         string
	Email        string
	PasswordHash []byte
}
