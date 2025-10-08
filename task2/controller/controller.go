package controller

import (
	"golang-training/task2/model"
	"golang-training/task2/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	UserService *service.UserService
}

func NewController(s *service.UserService) *Controller {
	return &Controller{UserService: s}
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Create a new user with name, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      model.User  true  "User Info"
// @Success      201   {object}  model.UserResponse
// @Failure      400   {object}  string
// @Router       /register [post]
func (ctrl *Controller) RegisterUser(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.UserService.RegisterUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Login with email and password to receive a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body  map[string]string  true  "User credentials"
// @Success      200  {object}  string
// @Failure      401  {object}  string
// @Router       /login [post]
func (ctrl *Controller) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.UserService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate the current JWT token
// @Tags         auth
// @Produce      json
// @Success      200  {object}  string
// @Failure      401  {object}  string
// @Router       /logout [post]
// @Security     ApiKeyAuth
func (ctrl *Controller) Logout(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	token, _ := c.Get("token")
	ctrl.UserService.LogoutUser(token.(string))

	c.JSON(http.StatusOK, gin.H{"message": "logged out", "user": user})
}

// GetCurrentUser godoc
// @Summary      Get current authenticated user
// @Description  Retrieve the details of the user based on the provided JWT token
// @Tags         users
// @Produce      json
// @Success      200  {object}  model.UserResponse
// @Failure      401  {object}  string
// @Router       /me [get]
// @Security     ApiKeyAuth
func (ctrl *Controller) GetCurrentUser(c *gin.Context) {
	existing, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	user, ok := existing.(model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user from context"})
		return
	}

	c.JSON(http.StatusOK, model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Retrieve all registered users (requires authentication)
// @Tags         users
// @Produce      json
// @Success      200  {array}   model.UserResponse
// @Failure      401  {object}  string
// @Router       /users [get]
// @Security     ApiKeyAuth
func (ctrl *Controller) GetAllUsers(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	users, err := ctrl.UserService.GetAllUsersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var safeUsers []model.UserResponse
	for _, user := range users {
		safeUsers = append(safeUsers, model.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	c.JSON(http.StatusOK, safeUsers)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieve a user’s details using their unique ID (requires authentication)
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  model.UserResponse
// @Failure      401  {object}  string
// @Failure      404  {object}  string
// @Router       /users/{id} [get]
// @Security     ApiKeyAuth
func (ctrl *Controller) GetUserByID(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id := c.Param("id")
	user, err := ctrl.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to fetch user"})
		return
	}

	c.JSON(http.StatusCreated, model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
