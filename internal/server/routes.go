package server

import (
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/auth"
	"github.com/davelawson/grim/lobby"
	"github.com/davelawson/grim/match"
	"github.com/davelawson/grim/model"
	"github.com/davelawson/grim/user"
	"github.com/davelawson/grim/util"
	"net/url"

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
	router.Handle(api.GetUserByEmail.Method, api.GetUserByEmail.Route, createAuthedHandler(authService, router, controller.GetUserByEmail))
}

func AddLobbyRoutes(authService authService, router *gin.Engine, controller *lobby.Controller) {
	router.Handle(api.ListPublicLobbies.Method, api.ListPublicLobbies.Route, createAuthedHandler(authService, router, controller.ListPublic))
	router.Handle(api.JoinLobby.Method, api.JoinLobby.Route, createAuthedHandler(authService, router, controller.Join))
	router.Handle(api.ListLobbies.Method, api.ListLobbies.Route, createAuthedHandler(authService, router, controller.List))
	router.Handle(api.CreateLobby.Method, api.CreateLobby.Route, createAuthedHandler(authService, router, controller.CreateLobby))
	router.Handle(api.DeleteLobby.Method, api.DeleteLobby.Route, createAuthedHandler(authService, router, controller.DeleteLobby))
	router.Handle(api.GetLobby.Method, api.GetLobby.Route, createAuthedHandler(authService, router, controller.GetLobby))
	router.Handle(api.UpdateLobby.Method, api.UpdateLobby.Route, createAuthedHandler(authService, router, controller.UpdateLobby))
	router.Handle(api.AddUserToLobby.Method, api.AddUserToLobby.Route, createAuthedHandler(authService, router, controller.AddUserToLobby))
	router.Handle(api.LeaveLobby.Method, api.LeaveLobby.Route, createAuthedHandler(authService, router, controller.Leave))
	router.Handle(api.RemoveUserFromLobby.Method, api.RemoveUserFromLobby.Route, createAuthedHandler(authService, router, controller.RemoveUserFromLobby))
	router.Handle(api.SetReady.Method, api.SetReady.Route, createAuthedHandler(authService, router, controller.SetReady))
	router.Handle(api.LaunchMatch.Method, api.LaunchMatch.Route, createAuthedHandler(authService, router, controller.Launch))
}

func AddMatchRoutes(authService authService, router *gin.Engine, controller *match.Controller) {
	// Match parameters include instance and choice IDs containing literal slashes.
	router.UseRawPath = true
	router.UnescapePathValues = false
	router.Handle(api.PlayCard.Method, api.PlayCard.Route, createAuthedHandler(authService, router, controller.PlayCard))
	router.Handle(api.Attack.Method, api.Attack.Route, createAuthedHandler(authService, router, controller.Attack))
	router.Handle(api.AttemptChallenge.Method, api.AttemptChallenge.Route, createAuthedHandler(authService, router, controller.AttemptChallenge))
	router.Handle(api.ActivateAbility.Method, api.ActivateAbility.Route, createAuthedHandler(authService, router, controller.ActivateAbility))
	router.Handle(api.ResearchElement.Method, api.ResearchElement.Route, createAuthedHandler(authService, router, controller.ResearchElement))
	router.Handle(api.ResearchSpell.Method, api.ResearchSpell.Route, createAuthedHandler(authService, router, controller.ResearchSpell))
	router.Handle(api.HealWithFocus.Method, api.HealWithFocus.Route, createAuthedHandler(authService, router, controller.HealWithFocus))
	router.Handle(api.EndTurn.Method, api.EndTurn.Route, createAuthedHandler(authService, router, controller.EndTurn))
	router.Handle(api.Maintenance.Method, api.Maintenance.Route, createAuthedHandler(authService, router, controller.Maintenance))
	router.Handle(api.Kneel.Method, api.Kneel.Route, createAuthedHandler(authService, router, controller.Kneel))
	router.Handle(api.ResolveEffect.Method, api.ResolveEffect.Route, createAuthedHandler(authService, router, controller.ResolveEffect))
	router.Handle(api.SkipEffect.Method, api.SkipEffect.Route, createAuthedHandler(authService, router, controller.SkipEffect))
	router.Handle(api.UndertakeHunt.Method, api.UndertakeHunt.Route, createAuthedHandler(authService, router, controller.UndertakeHunt))
	router.Handle(api.ChooseResearch.Method, api.ChooseResearch.Route, createAuthedHandler(authService, router, controller.ChooseResearch))
	router.Handle(api.EndDiscard.Method, api.EndDiscard.Route, createAuthedHandler(authService, router, controller.EndDiscard))
	router.Handle(api.HandDiscard.Method, api.HandDiscard.Route, createAuthedHandler(authService, router, controller.HandDiscard))
	router.Handle(api.AssetLoss.Method, api.AssetLoss.Route, createAuthedHandler(authService, router, controller.AssetLoss))
	router.Handle(api.Scry.Method, api.Scry.Route, createAuthedHandler(authService, router, controller.Scry))
	router.Handle(api.React.Method, api.React.Route, createAuthedHandler(authService, router, controller.React))
	router.Handle(api.Defend.Method, api.Defend.Route, createAuthedHandler(authService, router, controller.Defend))
	router.Handle(api.Pass.Method, api.Pass.Route, createAuthedHandler(authService, router, controller.Pass))
	router.Handle(api.ListMatches.Method, api.ListMatches.Route, createAuthedHandler(authService, router, controller.List))
	router.Handle(api.GetMatch.Method, api.GetMatch.Route, createAuthedHandler(authService, router, controller.Get))
	router.Handle(api.GetDraft.Method, api.GetDraft.Route, createAuthedHandler(authService, router, controller.GetDraft))
	router.Handle(api.SubmitDraft.Method, api.SubmitDraft.Route, createAuthedHandler(authService, router, controller.SubmitDraft))
	router.Handle(api.GetCatalogue.Method, api.GetCatalogue.Route, createAuthedHandler(authService, router, controller.Catalogue))
	router.Handle(api.EndMatch.Method, api.EndMatch.Route, createAuthedHandler(authService, router, controller.End))
}

// createAuthedHandler authenticates protected routes and returns JSON errors.
func createAuthedHandler(authService authService, router *gin.Engine, handler func(*gin.Context)) func(*gin.Context) {
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
		// Decode path segments once; PathUnescape preserves literal plus signs.
		if router.UseRawPath && !router.UnescapePathValues && c.Request.URL.RawPath != "" {
			for i := range c.Params {
				value, err := url.PathUnescape(c.Params[i].Value)
				if err != nil {
					util.WriteAPIError(c, util.ErrInvalidRequest)
					return
				}
				c.Params[i].Value = value
			}
		}
		c.Set("reqUser", reqUser)
		handler(c)
	}
}
