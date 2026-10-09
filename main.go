package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"main/api"
	"main/auth"
	"main/docs"
	"main/lobby"
	"main/match"
	"main/user"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type config struct {
	DbLocation string `json:"grim_db"`
	SslFolder  string `json:"grim_ssl"`
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

	db, err := sql.Open("sqlite3", config.DbLocation+"?_foreign_keys=on&_txlock=immediate&_busy_timeout=5000")
	if err != nil {
		fmt.Println("Error opening database at ", config.DbLocation, ": ", err)
		return
	}

	router := gin.Default()

	userRepo := user.NewUserRepo(db)
	lobbyRepo := lobby.NewLobbyRepo(db)

	authService := auth.NewService(userRepo)
	authFacade := auth.NewServiceFacade(authService, db)
	authController := auth.NewController(authFacade)
	api.AddAuthRoutes(router, authController)

	userService := user.NewService(userRepo)
	userFacade := user.NewServiceFacade(userService)
	userController := user.NewController(userFacade)
	api.AddUserRoutes(authFacade, router, userController)

	matchService := match.NewService(match.NewRepo(db), userRepo)
	matchController := match.NewController(match.NewServiceFacade(matchService))
	api.AddMatchRoutes(authFacade, router, matchController)
	lobbyService := lobby.NewService(lobbyRepo, userRepo, matchService)
	lobbyFacade := lobby.NewServiceFacade(lobbyService)
	lobbyController := lobby.NewController(lobbyFacade)
	api.AddLobbyRoutes(authFacade, router, lobbyController)

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
