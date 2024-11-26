package controller

import (
	"fmt"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	userService *service.UserService
}

type getUserRequest struct {
	Email string
}

type createUserRequest struct {
	Email        string
	Name         string
	PasswordHash string
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (us *UserController) GetUserByEmail(c *gin.Context) {
	// Do we really want the controllers to be coupled to Gin-Gonic?  Maybe we can find a way to handle this more generically...
	// Maybe creating a base class?  Maybe a delegate?  Maybe a wrapper for the context?
	// TODO: handle errors
	req := getUserRequest{}
	fmt.Println("GetUserByEmail(): {}", req)
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.Error(reqErr)
		return
	}

	user, err := us.userService.GetUserByEmail(req.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "{}", err)
		return
	}
	if user == nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (us *UserController) CreateUser(c *gin.Context) {
	req := createUserRequest{}
	fmt.Println("CreateUser(): {}", req)
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		c.Error(reqErr)
		return
	}

	err := us.userService.CreateUser(req.Email, req.Name, req.PasswordHash)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong.  Unable to create user. {}", err)
	}

	c.Status(http.StatusOK)
}
