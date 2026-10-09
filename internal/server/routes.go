package server

import (
	"main/api"
	"main/auth"
	"main/lobby"
	"main/match"
	"main/model"
	"main/user"
	"main/util"

	"github.com/gin-gonic/gin"
)

type authService interface {
	VerifyBearerToken(token string) (*model.User, error)
}

func AddAuthRoutes(router *gin.Engine, controller *auth.Controller) {
	router.Handle(api.LoginMethod, api.LoginRoute, controller.Login)
}

func AddUserRoutes(authService authService, router *gin.Engine, controller *user.Controller) {
	router.Handle(api.CreateUserMethod, api.CreateUserRoute, controller.CreateUser)
	router.Handle(api.GetUserByEmailMethod, api.GetUserByEmailRoute, createAuthedHandler(authService, controller.GetUserByEmail))
}

func AddLobbyRoutes(authService authService, router *gin.Engine, controller *lobby.Controller) {
	router.Handle(api.CreateLobbyMethod, api.CreateLobbyRoute, createAuthedHandler(authService, controller.CreateLobby))
	router.Handle(api.DeleteLobbyMethod, api.DeleteLobbyRoute, createAuthedHandler(authService, controller.DeleteLobby))
	router.Handle(api.GetLobbyMethod, api.GetLobbyRoute, createAuthedHandler(authService, controller.GetLobby))
	router.Handle(api.UpdateLobbyMethod, api.UpdateLobbyRoute, createAuthedHandler(authService, controller.UpdateLobby))
	router.Handle(api.AddUserToLobbyMethod, api.AddUserToLobbyRoute, createAuthedHandler(authService, controller.AddUserToLobby))
	router.Handle(api.RemoveUserFromLobbyMethod, api.RemoveUserFromLobbyRoute, createAuthedHandler(authService, controller.RemoveUserFromLobby))
	router.Handle(api.SetReadyMethod, api.SetReadyRoute, createAuthedHandler(authService, controller.SetReady))
	router.Handle(api.LaunchMatchMethod, api.LaunchMatchRoute, createAuthedHandler(authService, controller.Launch))
}

func AddMatchRoutes(authService authService, router *gin.Engine, controller *match.Controller) {
	router.Handle(api.GetMatchMethod, api.GetMatchRoute, createAuthedHandler(authService, controller.Get))
	router.Handle(api.SubmitMatchCommandMethod, api.SubmitMatchCommandRoute, createAuthedHandler(authService, controller.Command))
	router.Handle(api.GetCatalogueMethod, api.GetCatalogueRoute, createAuthedHandler(authService, controller.Catalogue))
	router.Handle(api.EndMatchMethod, api.EndMatchRoute, createAuthedHandler(authService, controller.End))
}

// createAuthedHandler authenticates protected routes and returns JSON errors.
func createAuthedHandler(authService authService, handler func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		reqUser, err := authService.VerifyBearerToken(util.GetBearerToken(c))
		if err != nil {
			util.WriteAPIError(c, err)
			return
		}
		if reqUser == nil {
			util.WriteAPIError(c, util.ErrUnauthorized)
			return
		}
		c.Set("reqUser", reqUser)
		handler(c)
	}
}
