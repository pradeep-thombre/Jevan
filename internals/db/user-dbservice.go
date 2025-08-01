package db

import (
	"Jevan/commons"
	"Jevan/commons/appdb"
	"Jevan/commons/apploggers"
	"Jevan/configs"
	"Jevan/internals/models"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type udbservice struct {
	ucollection appdb.DatabaseCollection
	dcollection appdb.DatabaseCollection
}

type UserDbService interface {
	GetUserById(ctx context.Context, id string) (*models.User, error)
	DeleteUserById(ctx context.Context, id string) error
	GetUsers(ctx context.Context) ([]models.User, error)
	CreateUserProfile(ctx context.Context, user *models.User) (string, error)
	UpdateUser(ctx context.Context, user *models.User, userId string) error
	RegisterUser(ctx context.Context, user *models.UserDetails) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error)
	UpdateUserRole(ctx context.Context, userID string, newRole string) error
}

func NewUserDbService(dbclient appdb.DatabaseClient) UserDbService {
	return &udbservice{
		ucollection: dbclient.Collection(configs.MONGO_USERS_COLLECTION),
		dcollection: dbclient.Collection(configs.MONGO_USERDETAILS_COLLECTION),
	}
}

func (u *udbservice) GetUserById(ctx context.Context, userId string) (*models.User, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing GetUserById, Id: %s", userId)
	// get object id from userid string
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid userid provided, userId: %s", userId)
	}
	var user *models.User
	var filter = bson.M{"_id": id}
	dbError := u.dcollection.FindOne(ctx, filter, &user)
	if dbError != nil {
		logger.Error(dbError)
		return nil, dbError
	}
	logger.Infof("Executed GetUserById, user: %s", commons.PrintStruct(user))
	return user, nil
}

func (u *udbservice) DeleteUserById(ctx context.Context, userId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing DeleteUserById, Id: %s", userId)
	// get object id from userid string
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("cannot delete user, invalid userid provided, userId: %s", userId)
	}
	var filter = bson.M{"_id": id}
	_, dbError := u.dcollection.DeleteOne(ctx, filter)
	if dbError != nil {
		logger.Error(dbError)
		return dbError
	}
	logger.Infof("Executed DeleteUserById, Id: %s", userId)
	return nil
}

func (u *udbservice) GetUsers(ctx context.Context) ([]models.User, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing GetUsers")

	// create users payload to find data from db
	var users []models.User
	var filter = map[string]interface{}{}
	dbError := u.dcollection.Find(ctx, filter, &options.FindOptions{}, &users)
	if dbError != nil {
		logger.Error(dbError)
		return nil, dbError
	}
	logger.Infof("Executed GetUsers, users: %d", len(users))
	return users, nil
}

func (u *udbservice) CreateUserProfile(ctx context.Context, user *models.User) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing SaveUser...")

	// insert user in db
	result, dbError := u.dcollection.InsertOne(ctx, user)
	if dbError != nil {
		logger.Error(dbError)
		return "", dbError
	}

	// Extract the inserted ID from the result
	id := result.InsertedID.(primitive.ObjectID).Hex()
	logger.Infof("Executed SaveUser, userid: %s", commons.PrintStruct(user))
	return id, nil
}

func (u *udbservice) UpdateUser(ctx context.Context, user *models.User, userId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing UpdateUser...")
	// get object id from userid string
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("cannot delete user, invalid userid provided, userId: %s", userId)
	}
	var filter = bson.M{"_id": id}
	update := bson.M{"$set": user} // Correct update document
	// update user in db
	_, dbError := u.dcollection.UpdateOne(ctx, filter, update)
	if dbError != nil {
		logger.Error(dbError)
		return dbError
	}

	logger.Infof("Executed UpdateUser, userid: %s", commons.PrintStruct(user))
	return nil
}

func (u *udbservice) RegisterUser(ctx context.Context, user *models.UserDetails) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Creating user: %s", user.Email)

	// check if user exists
	var existing models.UserDetails
	err := u.ucollection.FindOne(ctx, bson.M{"email": user.Email}, &existing)
	if err == nil {
		return "", errors.New("user already exists")
	}

	result, err := u.ucollection.InsertOne(ctx, user)
	if err != nil {
		logger.Error("Failed to create user: ", err)
		return "", err
	}
	userId := result.InsertedID.(primitive.ObjectID).Hex()
	logger.Infof("User created successfully with ID: %s", userId)
	return userId, nil
}

func (u *udbservice) GetUserByEmail(ctx context.Context, email string) (*models.UserDetails, error) {
	var user models.UserDetails
	err := u.ucollection.FindOne(ctx, bson.M{"email": email}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *udbservice) UpdateUserRole(ctx context.Context, userID string, newRole string) error {
	objId, err := primitive.ObjectIDFromHex(userID)
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Updating user role for ID: %s to %s", userID, newRole)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"role": newRole}}
	_, err = u.ucollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	return err
}
