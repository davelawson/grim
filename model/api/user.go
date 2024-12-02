package api

import "main/model"

type User struct {
	Id    string
	Name  string
	Email string
}

type GetUserRequest struct {
	Email string
}

type GetUserResponse struct {
	User User
}

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
}

func NewUser(user *model.User) User {
	return User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}
