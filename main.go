package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"main/auth"
	"main/docs"
	"main/lobby"
	"main/model"
	"main/repo"
	"main/user"
	"main/util"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type config struct {
	DbLocation string `json:"GRIM_DB"`
	SslFolder  string `json:"GRIM_SSL"`
}

type authService interface {
	VerifyBearerToken(token string) (*model.User, error)
}

// @SecurityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config, configErr := getConfig()
	if configErr != nil {
		fmt.Println("Config contains errors.  Aborting launch.", configErr)
		return
	}

	db, err := sql.Open("sqlite3", config.DbLocation+"?_foreign_keys=on")
	if err != nil {
		fmt.Println("Error opening database at ", config.DbLocation, ": ", err)
		return
	}

	router := gin.Default()

	userRepo := repo.NewUserRepo(db)
	lobbyRepo := lobby.NewLobbyRepo(db)

	authService := auth.NewService(userRepo)
	authController := auth.NewController(authService)
	addAuthRoutes(router, authController)

	userService := user.NewService(userRepo)
	userController := user.NewController(userService)
	addUserRoutes(authService, router, userController)

	lobbyService := lobby.NewService(lobbyRepo, userRepo)
	lobbyFacade := lobby.NewServiceFacade(lobbyService)
	lobbyController := lobby.NewController(lobbyFacade)
	addLobbyRoutes(authService, router, lobbyController)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.RunTLS("localhost:8080", config.SslFolder+"/grim.crt", config.SslFolder+"/grim.key")
}

func getConfig() (*config, error) {
	c := &config{
		DbLocation: os.Getenv("GRIM_DB"),
		SslFolder:  os.Getenv("GRIM_SSL"),
	}
	json, _ := json.MarshalIndent(*c, "", "  ")
	fmt.Println("config: ")
	fmt.Println(string(json))

	errorMessage := ""
	if c.DbLocation == "" {
		errorMessage += "GRIM_DB environment variable missing"
	}
	if c.SslFolder == "" {
		errorMessage += "GRIM_SSL environment variable missing"
	}

	if errorMessage == "" {
		return c, nil
	} else {
		return nil, errors.New(errorMessage)
	}
}

func addSwaggerInfo() {
	docs.SwaggerInfo.Title = "Grimoire Backend API"
	docs.SwaggerInfo.Description = "This is the backend RESTful server for the Grimoire game."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"https"}
}

func addAuthRoutes(router *gin.Engine, controller *auth.Controller) {
	group := router.Group("/login")
	group.POST("", controller.Login)
}

func addUserRoutes(authService authService, router *gin.Engine, controller *user.Controller) {
	group := router.Group("/user")
	group.POST("", controller.CreateUser)

	group = router.Group("/user/getbyemail")
	group.POST("", createAuthedHandler(authService, controller.GetUserByEmail))
}

func addLobbyRoutes(authService authService, router *gin.Engine, controller *lobby.Controller) {
	group := router.Group("/lobby")
	group.POST("", createAuthedHandler(authService, controller.CreateLobby))
	group.DELETE(":id", createAuthedHandler(authService, controller.DeleteLobby))
	group.GET(":id", createAuthedHandler(authService, controller.GetLobby))
	group.PUT(":id", createAuthedHandler(authService, controller.UpdateLobby))
	group.POST(":id/user", createAuthedHandler(authService, controller.AddUserToLobby))
	group.DELETE(":id/user/:user_id", createAuthedHandler(authService, controller.RemoveUserFromLobby))
}

func createAuthedHandler(authService authService, handler func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		reqUser, authErr := authService.VerifyBearerToken(util.GetBearerToken(c))
		if authErr != nil {
			c.String(http.StatusInternalServerError, "Invalid authentication token", authErr)
			return
		}
		if reqUser == nil {
			c.String(http.StatusUnauthorized, "Bad or missing authentication token")
			return
		}
		c.Set("reqUser", reqUser)
		handler(c)
	}
}
