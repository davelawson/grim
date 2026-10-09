package api

import (
	"encoding/json"
	"errors"
	"main/model"
	"strings"
)

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUserRequest struct {
	Email string `json:"email"`
}

type GetUserResponse struct {
	User User `json:"user"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// UnmarshalJSON enforces registration field casing without changing decoding
// for lookup requests or introducing the remaining registration validation.
func (req *CreateUserRequest) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	for field := range fields {
		for _, name := range []string{"email", "name", "password"} {
			if field != name && strings.EqualFold(field, name) {
				return errors.New("registration field names must be lowercase")
			}
		}
	}

	// The alias avoids invoking this method recursively.
	type request CreateUserRequest
	return json.Unmarshal(data, (*request)(req))
}

func NewUser(user *model.User) User {
	return User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}
