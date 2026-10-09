package user

import (
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/model"
)

func userResponse(user *model.User) api.User {
	return api.User{Id: user.Id, Name: user.Name, Email: user.Email}
}
