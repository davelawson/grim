package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"main/controller"
	"main/docs"
	"main/repo"
	"main/service"
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

func main() {
	config, configErr := getConfig()
	if configErr != nil {
		fmt.Println("Config contains errors.  Aborting launch.", configErr)
		return
	}
	db, err := sql.Open("sqlite3", config.DbLocation)
	if err != nil {
		fmt.Println("Error opening database at {}: {}", config.DbLocation, err)
		return
	}

	router := gin.Default()

	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	addUserRoutes(router, userController)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.RunTLS("localhost:8080", config.SslFolder+"/grim.crt", config.SslFolder+"/grim.key")
}

func getConfig() (*config, error) {
	c := &config{
		DbLocation: os.Getenv("GRIM_DB"),
		SslFolder:  os.Getenv("GRIM_SSL"),
	}
	json, _ := json.MarshalIndent(*c, "", "  ")
	fmt.Println("config:")
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
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func addUserRoutes(router *gin.Engine, controller *controller.UserController) {
	group := router.Group("/user")
	group.POST("", controller.CreateUser)

	group = router.Group("/user/getbyemail")
	group.POST("", controller.GetUserByEmail)
}
