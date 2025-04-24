package db

import (
	"Jevan/commons/appdb"
	"Jevan/commons/apploggers"
	"Jevan/configs"
	"Jevan/internals/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cDbService struct {
	ucollection appdb.DatabaseCollection
}

type CartDbService interface {
	GetCartById(ctx context.Context, cartId string) (*models.Cart, error)
	SaveCart(ctx context.Context, cart *models.Cart) error
	DeleteAllItemsFromCart(ctx context.Context, cartId string) error
}

func NewCartDbService(dbclient appdb.DatabaseClient) CartDbService {
	return &cDbService{
		ucollection: dbclient.Collection(configs.MONGO_CARTS_COLLECTION),
	}
}

// Get Cart by cartId
func (c *cDbService) GetCartById(ctx context.Context, cartId string) (*models.Cart, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing GetCartById, cartId: %s", cartId)

	var cart models.Cart
	cartObjId, oerror := primitive.ObjectIDFromHex(cartId)
	if oerror != nil {
		return nil, errors.New("error: invalid id provided")
	}
	filter := bson.M{"_id": cartObjId}
	dbError := c.ucollection.FindOne(ctx, filter, &cart)
	if dbError != nil {
		logger.Error(dbError)
		return nil, dbError
	}

	logger.Infof("Executed GetCartById, cartId: %s", cartId)
	return &cart, nil
}

// Delete all items from the cart
func (c *cDbService) DeleteAllItemsFromCart(ctx context.Context, cartId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing DeleteAllItemsFromCart, cartId: %s", cartId)
	cartObjId, oerror := primitive.ObjectIDFromHex(cartId)
	if oerror != nil {
		return errors.New("error: invalid id provided")
	}
	filter := bson.M{"_id": cartObjId}
	update := bson.M{"$set": bson.M{"items": []models.CartItem{}, "totalprice": 0}}

	_, dbError := c.ucollection.UpdateOne(ctx, filter, update)
	if dbError != nil {
		logger.Error(dbError)
		return dbError
	}

	logger.Infof("Executed DeleteAllItemsFromCart, cartId: %s", cartId)
	return nil
}

func (c *cDbService) SaveCart(ctx context.Context, cart *models.Cart) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing SaveCart")

	if cart.ID.IsZero() {
		logger.Error("Cart ID is required")
		return errors.New("cart ID is required")
	}

	filter := bson.M{"_id": cart.ID}
	var existing models.Cart

	err := c.ucollection.FindOne(ctx, filter, &existing)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			// No document found, insert as new
			_, insertErr := c.ucollection.InsertOne(ctx, cart)
			if insertErr != nil {
				logger.Error(insertErr)
				return insertErr
			}
			logger.Infof("Inserted new cart with ID: %s", cart.ID.Hex())
			return nil
		}
		logger.Error(err)
		return err
	}

	// Document exists, update it
	update := bson.M{"$set": cart}
	_, updateErr := c.ucollection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		logger.Error(updateErr)
		return updateErr
	}

	logger.Infof("Updated existing cart with ID: %s", cart.ID.Hex())
	return nil
}
