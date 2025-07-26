package apis

import (
	"Jevan/commons"
	"Jevan/commons/apploggers"
	"Jevan/internals/models"
	"Jevan/internals/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ucontroller struct {
	eservice services.UserService
}

func NewUserController(eservice services.UserService) ucontroller {
	return ucontroller{
		eservice: eservice,
	}
}

// @Tags User Management
// @Summary GetUserById
// @Description Gets user details by user id such as name, email, status etc.
// @Accept json
// @Produce json
// @Param id path string true "User id"
// @Success 200 {object} models.User
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /users/{id} [Get]
func (u *ucontroller) GetUserById(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	userId := c.Param("id")
	logger.Infof("Executing GetUserById, userId: %s", userId)
	if len(strings.TrimSpace(userId)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}
	user, serror := u.eservice.GetUserById(lcontext, userId)
	if serror != nil {
		logger.Error(serror)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(serror.Error(), nil))
	}
	logger.Infof("Executed GetUserById, userId:%s, user %s", userId, commons.PrintStruct(user))
	return c.JSON(http.StatusOK, user)
}

// @Tags User Management
// @Summary DeleteUserById
// @Description delete user details by user id
// @Accept json
// @Produce json
// @Param id path string true "User id"
// @Success 204
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /users/{id} [Delete]
func (u *ucontroller) DeleteUserById(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	userId := c.Param("id")
	logger.Infof("Executing DeleteUserById, userId: %s", userId)
	if len(strings.TrimSpace(userId)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}
	serror := u.eservice.DeleteUserById(lcontext, userId)
	if serror != nil {
		logger.Error(serror)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(serror.Error(), nil))
	}
	logger.Infof("Executed DeleteUserById, userId: %s", userId)
	return c.NoContent(http.StatusNoContent)
}

// @Tags User Management
// @Summary GetUsers
// @Description get details of all users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /users [Get]
func (u *ucontroller) GetUsers(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	logger.Info("Executing Get All Users")
	users, serror := u.eservice.GetUsers(lcontext)
	if serror != nil {
		logger.Error(serror)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(serror.Error(), nil))
	}
	logger.Infof("Executed GetUsers, users %s", commons.PrintStruct(users))
	return c.JSON(http.StatusOK, map[string]interface{}{
		"total": len(users),
		"users": users,
	})
}

// @Tags User Management
// @Summary CreateUser
// @Description Create a user with name, email, age, and is_Active status
// @Accept json
// @Produce json
// @Param payload body models.User true "User data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /users [post]
func (u *ucontroller) CreateUser(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	logger.Info("Executing CreateUser")
	var user *models.User
	err := c.Bind(&user)
	if err != nil || user == nil {
		logger.Error("invalid request payload")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("invalid request payload", nil))
	}

	if len(strings.TrimSpace(user.Name)) == 0 {
		logger.Error("'name' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'name' is required", nil))
	}

	if len(strings.TrimSpace(user.Email)) == 0 {
		logger.Error("'email' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'email' is required", nil))
	}
	Id, serror := u.eservice.CreateUser(lcontext, user)
	if serror != nil {
		logger.Error(serror)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(serror.Error(), nil))
	}
	logger.Info("Executed CreateUser")
	return c.JSON(http.StatusCreated, map[string]string{
		"id": Id,
	})
}

// @Tags User Management
// @Summary UpdateUser
// @Description update user details such as name, email, age, and is_Active status bu user id
// @Accept json
// @Produce json
// @Param payload body models.User true "User data"
// @Param id path string true "User Id"
// @Success 200
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /users/{id} [patch]
func (u *ucontroller) UpdateUser(c echo.Context) error {
	lcontext, logger := apploggers.GetLoggerFromEcho(c)
	userId := c.Param("id")
	logger.Info("Executing UpdateUser, userId: %s", userId)
	if len(strings.TrimSpace(userId)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}
	var user *models.User
	err := c.Bind(&user)
	if err != nil || user == nil {
		logger.Error("invalid request payload")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("invalid request payload", nil))
	}

	if len(strings.TrimSpace(user.Name)) == 0 {
		logger.Error("'name' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'name' is required", nil))
	}

	if len(strings.TrimSpace(user.Email)) == 0 {
		logger.Error("'email' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'email' is required", nil))
	}
	serror := u.eservice.UpdateUser(lcontext, user, userId)
	if serror != nil {
		logger.Error(serror)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse(serror.Error(), nil))
	}
	logger.Info("Executed UpdateUser, userId: %s", userId)
	return c.NoContent(http.StatusOK)
}
