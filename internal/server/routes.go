package server

import (
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/auth"
	"github.com/davelawson/grim/lobby"
	"github.com/davelawson/grim/match"
	"github.com/davelawson/grim/model"
	"github.com/davelawson/grim/user"
	"github.com/davelawson/grim/util"

	"github.com/gin-gonic/gin"
)

type authService interface {
	VerifyBearerToken(token string) (*model.User, error)
}

func AddAuthRoutes(router *gin.Engine, controller *auth.Controller) {
	router.Handle(api.Login.Method, api.Login.Route, controller.Login)
}

func AddUserRoutes(authService authService, router *gin.Engine, controller *user.Controller) {
	router.Handle(api.CreateUser.Method, api.CreateUser.Route, controller.CreateUser)
	router.Handle(api.GetUserByEmail.Method, api.GetUserByEmail.Route, createAuthedHandler(authService, controller.GetUserByEmail))
}

func AddLobbyRoutes(authService authService, router *gin.Engine, controller *lobby.Controller) {
	router.Handle(api.ListPublicLobbies.Method, api.ListPublicLobbies.Route, createAuthedHandler(authService, controller.ListPublic))
	router.Handle(api.JoinLobby.Method, api.JoinLobby.Route, createAuthedHandler(authService, controller.Join))
	router.Handle(api.ListLobbies.Method, api.ListLobbies.Route, createAuthedHandler(authService, controller.List))
	router.Handle(api.CreateLobby.Method, api.CreateLobby.Route, createAuthedHandler(authService, controller.CreateLobby))
	router.Handle(api.DeleteLobby.Method, api.DeleteLobby.Route, createAuthedHandler(authService, controller.DeleteLobby))
	router.Handle(api.GetLobby.Method, api.GetLobby.Route, createAuthedHandler(authService, controller.GetLobby))
	router.Handle(api.UpdateLobby.Method, api.UpdateLobby.Route, createAuthedHandler(authService, controller.UpdateLobby))
	router.Handle(api.AddUserToLobby.Method, api.AddUserToLobby.Route, createAuthedHandler(authService, controller.AddUserToLobby))
	router.Handle(api.LeaveLobby.Method, api.LeaveLobby.Route, createAuthedHandler(authService, controller.Leave))
	router.Handle(api.RemoveUserFromLobby.Method, api.RemoveUserFromLobby.Route, createAuthedHandler(authService, controller.RemoveUserFromLobby))
	router.Handle(api.SetReady.Method, api.SetReady.Route, createAuthedHandler(authService, controller.SetReady))
	router.Handle(api.LaunchMatch.Method, api.LaunchMatch.Route, createAuthedHandler(authService, controller.Launch))
}

func AddMatchRoutes(authService authService, router *gin.Engine, controller *match.Controller) {
	router.Handle(api.ListMatches.Method, api.ListMatches.Route, createAuthedHandler(authService, controller.List))
	router.Handle(api.GetMatch.Method, api.GetMatch.Route, createAuthedHandler(authService, controller.Get))
	router.Handle(api.GetDraft.Method, api.GetDraft.Route, createAuthedHandler(authService, controller.GetDraft))
	router.Handle(api.SubmitDraft.Method, api.SubmitDraft.Route, createAuthedHandler(authService, controller.SubmitDraft))
	router.Handle(api.GetCatalogue.Method, api.GetCatalogue.Route, createAuthedHandler(authService, controller.Catalogue))
	router.Handle(api.EndMatch.Method, api.EndMatch.Route, createAuthedHandler(authService, controller.End))
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
