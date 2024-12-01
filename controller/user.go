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
	Email    string
	Name     string
	Password string
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetUserByEmail godoc
//
//	@Summary		Get user by email
//	@Description	Lookup a specific user by email
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		getUserRequest	true	"Request Object"
//	@Success		200		{object}	model.User
//	@Router			/user/getbyemail [post]
func (us *UserController) GetUserByEmail(c *gin.Context) {
	// TODO: find a re-usable way to translate the context into a typed request
	req := getUserRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("GetUserByEmail(): ", req)
	if reqErr != nil {
		// TODO: Make this suck less
		c.Error(reqErr)
		return
	}

	user, err := us.userService.GetUserByEmail(req.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Woops", err)
		return
	}
	if user == nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
//
//	@Summary		Create user
//	@Description	Create a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body	createUserRequest	true	"Request Object"
//	@Success		200
//	@Router			/user [post]
func (us *UserController) CreateUser(c *gin.Context) {
	req := createUserRequest{}
	fmt.Println("CreateUser(): {}", req)
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		c.Error(reqErr)
		return
	}

	err := us.userService.CreateUser(req.Email, req.Name, req.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong.  Unable to create user. {}", err)
	}

	c.Status(http.StatusOK)
}
