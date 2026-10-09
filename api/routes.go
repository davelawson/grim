package api

import (
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
	group := router.Group("/login")
	group.POST("", controller.Login)
}

func AddUserRoutes(authService authService, router *gin.Engine, controller *user.Controller) {
	group := router.Group("/user")
	group.POST("", controller.CreateUser)

	group = router.Group("/user/getbyemail")
	group.POST("", createAuthedHandler(authService, controller.GetUserByEmail))
}

func AddLobbyRoutes(authService authService, router *gin.Engine, controller *lobby.Controller) {
	group := router.Group("/lobby")
	group.POST("", createAuthedHandler(authService, controller.CreateLobby))
	group.DELETE(":id", createAuthedHandler(authService, controller.DeleteLobby))
	group.GET(":id", createAuthedHandler(authService, controller.GetLobby))
	group.PUT(":id", createAuthedHandler(authService, controller.UpdateLobby))
	group.POST(":id/user", createAuthedHandler(authService, controller.AddUserToLobby))
	group.DELETE(":id/user/:user_id", createAuthedHandler(authService, controller.RemoveUserFromLobby))
	group.PUT(":id/ready", createAuthedHandler(authService, controller.SetReady))
	group.POST(":id/launch", createAuthedHandler(authService, controller.Launch))
}

func AddMatchRoutes(authService authService, router *gin.Engine, controller *match.Controller) {
	router.GET("/match/:id", createAuthedHandler(authService, controller.Get))
	router.POST("/match/:id/commands", createAuthedHandler(authService, controller.Command))
	router.GET("/match/:id/catalogue", createAuthedHandler(authService, controller.Catalogue))
	router.POST("/admin/match/:id/end", createAuthedHandler(authService, controller.End))
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
