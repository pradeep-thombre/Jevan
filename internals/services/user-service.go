package services

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/db"
	"Jevan/internals/models"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserById(context context.Context, userId string) (*models.User, error)
	DeleteUserById(context context.Context, userId string) error
	GetUsers(context context.Context) ([]models.User, error)
	CreateUserProfile(context context.Context, user *models.User) (string, error)
	UpdateUser(context context.Context, user *models.User, userId string) error
	RegisterUser(ctx context.Context, email, password string) (string, error)
	AuthenticateUser(ctx context.Context, email, password string) (*models.UserDetails, bool, error)
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

func (e *userService) GetUsers(context context.Context) ([]models.User, error) {
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

	user.CartId = primitive.NewObjectID().Hex()
	userId, dberror := e.dbservice.CreateUserProfile(context, user)
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

	dberror := e.dbservice.UpdateUser(context, user, userId)
	if dberror != nil {
		logger.Error(dberror)
		return dberror
	}
	logger.Infof("Executed UpdateUser, userId: %v", userId)
	return nil
}

func (s *userService) RegisterUser(ctx context.Context, email, password string) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Password hashing failed: ", err)
		return "", err
	}

	user := &models.UserDetails{
		Email:    email,
		Password: string(hashed),
		Role:     "user",
	}

	id, err := s.dbservice.RegisterUser(ctx, user)
	if err != nil {
		logger.Error("Failed to register user: ", err)
		return "", err
	}

	logger.Info("User registered successfully, Id: ", id)
	return id, nil
}

func (s *userService) AuthenticateUser(ctx context.Context, email, password string) (*models.UserDetails, bool, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Authenticating user: %s", email)

	user, err := s.dbservice.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, false, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Error("Password mismatch")
		return nil, false, errors.New("invalid password")
	}

	logger.Info("User authenticated successfully: ", email)
	return user, true, nil
}

func (s *userService) UpdateUserRole(ctx context.Context, userID string, newRole string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Updating role for user ID: %s to %s", userID, newRole)
	if newRole != "admin" && newRole != "user" {
		return fmt.Errorf("invalid role: %s", newRole)
	}

	return s.dbservice.UpdateUserRole(ctx, userID, newRole)
}
