package apis

import (
	"Jevan/commons"
	"Jevan/commons/apploggers"
	"Jevan/configs"
	"Jevan/internals/models"
	"Jevan/internals/services"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	userService services.UserService
}

func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{userService: userService}
}

// @Summary Register User
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.UserDetails true "User registration data"
// @Success 201 {string} string "Registered successfully"
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /register [post]
func (ac *AuthController) Register(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	logger.Info("Received registration request")
	var user models.UserDetails
	if err := c.Bind(&user); err != nil {
		logger.Error("Invalid request body: ", err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Invalid request body", nil))
	}

	if errs := commons.ValidateStruct(user); errs != nil {
		logger.Error("Validation error: ", errs)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Validation error: "+errs.Error(), nil))
	}

	if err := ac.userService.RegisterUser(lcontext, user.Email, user.Password); err != nil {
		logger.Error("Registration failed: ", err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Registration failed: "+err.Error(), nil))
	}

	userInfo := &models.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	// Create user profile in the database
	logger.Info("Creating user profile for: ", user.Email)
	id, err := ac.userService.CreateUserProfile(lcontext, userInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Failed to create user profile: "+err.Error(), nil))
	}

	logger.Info("User registered successfully with ID: ", id)
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Registered successfully",
		"id":      id,
	})

}

// @Summary Login User
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.UserDetails true "User credentials"
// @Success 200 {object} models.UserDetails
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /login [post]
func (ac *AuthController) Login(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	var creds models.UserDetails
	if err := c.Bind(&creds); err != nil {
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Invalid request body", nil))
	}
	logger.Info("Received login request for user: ", creds.Email)

	if errs := commons.ValidateStruct(creds); errs != nil {
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Validation error: "+errs.Error(), nil))
	}

	role, ok, err := ac.userService.AuthenticateUser(lcontext, creds.Email, creds.Password)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, commons.ApiErrorResponse("Invalid credentials", nil))
	}

	claims := jwt.MapClaims{
		"email": creds.Email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(configs.AppConfig.JwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Token generation failed", nil))
	}
	logger.Info("User logged in successfully: ", creds.Email)
	return c.JSON(http.StatusOK, models.UserLoginResponse{Token: signed})
}

// UpdateUserRole godoc
// @Summary Update user role (admin only)
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body models.UpdateUserRoleRequest true "New role (admin or user)"
// @Success 200 {object} map[string]string "Role updated successfully"
// @Failure 400 {object} map[string]string "Invalid request or validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /admin/users/{id}/role [put]
func (ac *AuthController) UpdateUserRole(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	id := c.Param("id")

	logger.Info("Received request to update role for user ID: ", id)
	var body models.UpdateUserRoleRequest

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	if err := commons.ValidateStruct(body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err := ac.userService.UpdateUserRole(lcontext, id, body.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Role updated successfully"})
}
