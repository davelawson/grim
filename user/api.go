package user

import (
	"main/api"
	"main/model"
)

func userResponse(user *model.User) api.User {
	return api.User{Id: user.Id, Name: user.Name, Email: user.Email}
}
