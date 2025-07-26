package services

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/db"
	dbmodel "Jevan/internals/db/models"
	"Jevan/internals/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserById(context context.Context, userId string) (*models.User, error)
	DeleteUserById(context context.Context, userId string) error
	GetUsers(context context.Context) ([]*models.User, error)
	CreateUserProfile(context context.Context, user *models.User) (string, error)
	UpdateUser(context context.Context, user *models.User, userId string) error
	RegisterUser(ctx context.Context, email, password string) error
	AuthenticateUser(ctx context.Context, email, password string) (string, bool, error)
	UpdateUserRole(ctx context.Context, userID, newRole string) error
}

type userService struct {
	dbservice db.UserDbService
}

func NewUserService(dbservice db.UserDbService) UserService {
	return &userService{
		dbservice: dbservice,
	}
}

func (e *userService) GetUserById(context context.Context, userId string) (*models.User, error) {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	logger.Infof("Executing GetUserById, userId: %s", userId)
	user, dberror := e.dbservice.GetUserById(context, userId)
	if dberror != nil {
		logger.Error(dberror)
		return nil, dberror
	}
	logger.Infof("Executed GetUserById, userId: %s", userId)
	return user, nil
}

func (e *userService) DeleteUserById(context context.Context, userId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	logger.Infof("Executing DeleteUserById, userId: %s", userId)
	dberror := e.dbservice.DeleteUserById(context, userId)
	if dberror != nil {
		logger.Error(dberror)
		return dberror
	}
	logger.Infof("Executed DeleteUserById, userId: %s", userId)
	return nil
}

func (e *userService) GetUsers(context context.Context) ([]*models.User, error) {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	logger.Infof("Executing GetUsers...")
	users, dberror := e.dbservice.GetUsers(context)
	if dberror != nil {
		logger.Error(dberror)
		return nil, dberror
	}
	logger.Infof("Executed GetUsers, users: %d", len(users))
	return users, nil
}

func (e *userService) CreateUserProfile(context context.Context, user *models.User) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	logger.Infof("Executing CreateUserProfile...")
	var userSchema *dbmodel.UserSchema
	pbyes, _ := json.Marshal(user)
	uerror := json.Unmarshal(pbyes, &userSchema)
	if uerror != nil {
		logger.Error(uerror.Error())
		return "", uerror
	}
	userId, dberror := e.dbservice.CreateUserProfile(context, userSchema)
	if dberror != nil {
		logger.Error(dberror)
		return "", dberror
	}
	logger.Infof("Executed CreateUserProfile, userId: %v", userId)
	return userId, nil
}

func (e *userService) UpdateUser(context context.Context, user *models.User, userId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(context)
	logger.Infof("Executing UpdateUser...")
	var userSchema *dbmodel.UserSchema
	pbyes, _ := json.Marshal(user)
	uerror := json.Unmarshal(pbyes, &userSchema)
	if uerror != nil {
		logger.Error(uerror.Error())
		return uerror
	}
	dberror := e.dbservice.UpdateUser(context, userSchema, userId)
	if dberror != nil {
		logger.Error(dberror)
		return dberror
	}
	logger.Infof("Executed UpdateUser, userId: %v", userId)
	return nil
}

func (s *userService) RegisterUser(ctx context.Context, email, password string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Password hashing failed: ", err)
		return err
	}

	user := &models.UserDetails{
		Email:    email,
		Password: string(hashed),
		Role:     "user",
	}

	err = s.dbservice.RegisterUser(ctx, user)
	if err != nil {
		logger.Error("Failed to register user: ", err)
		return err
	}

	logger.Info("User registered successfully")
	return nil
}

func (s *userService) AuthenticateUser(ctx context.Context, email, password string) (string, bool, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Authenticating user: %s", email)

	user, err := s.dbservice.GetUserByEmail(ctx, email)
	if err != nil {
		return "", false, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Error("Password mismatch")
		return "", false, errors.New("invalid password")
	}

	logger.Info("User authenticated successfully: ", email)
	return user.Role, true, nil
}

func (s *userService) UpdateUserRole(ctx context.Context, userID string, newRole string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Updating role for user ID: %s to %s", userID, newRole)
	if newRole != "admin" && newRole != "user" {
		return fmt.Errorf("invalid role: %s", newRole)
	}

	return s.dbservice.UpdateUserRole(ctx, userID, newRole)
}
