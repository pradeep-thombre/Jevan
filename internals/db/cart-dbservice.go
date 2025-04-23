package db

import (
	"Jevan/commons/appdb"
	"Jevan/commons/apploggers"
	"Jevan/configs"
	"Jevan/internals/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type cDbService struct {
	ucollection appdb.DatabaseCollection
}

type CartDbService interface {
	AddToCart(ctx context.Context, cartId string, item *models.CartItem) error
	GetCartById(ctx context.Context, cartId string) (*models.Cart, error)
	DeleteItemFromCart(ctx context.Context, cartId string, itemId string) error
	DeleteAllItemsFromCart(ctx context.Context, cartId string) error
}

func NewCartDbService(dbclient appdb.DatabaseClient) CartDbService {
	return &cDbService{
		ucollection: dbclient.Collection(configs.MONGO_CARTS_COLLECTION),
	}
}

// Add item to a cart
func (c *cDbService) AddToCart(ctx context.Context, cartId string, item *models.CartItem) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing AddToCart, cartId: %s", cartId)

	// Find the cart by cartId
	var cart models.Cart
	filter := bson.M{"cartId": cartId}
	dbError := c.ucollection.FindOne(ctx, filter, &cart)
	if dbError != nil {
		if dbError.Error() == "mongo: no documents in result" {
			// Cart not found, create a new one
			cart.ID = cartId
			cart.Items = []models.CartItem{*item}
			_, insertError := c.ucollection.InsertOne(ctx, &cart)
			if insertError != nil {
				logger.Error(insertError)
				return insertError
			}
			return nil
		}
		logger.Error(dbError)
		return dbError
	}

	// If cart exists, add the item to the cart
	cart.Items = append(cart.Items, *item)
	update := bson.M{"$set": bson.M{"items": cart.Items}}
	_, updateError := c.ucollection.UpdateOne(ctx, filter, update)
	if updateError != nil {
		logger.Error(updateError)
		return updateError
	}

	logger.Infof("Executed AddToCart, cartId: %s", cartId)
	return nil
}

// Get Cart by cartId
func (c *cDbService) GetCartById(ctx context.Context, cartId string) (*models.Cart, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing GetCartById, cartId: %s", cartId)

	var cart models.Cart
	filter := bson.M{"cartId": cartId}
	dbError := c.ucollection.FindOne(ctx, filter, &cart)
	if dbError != nil {
		logger.Error(dbError)
		return nil, dbError
	}

	logger.Infof("Executed GetCartById, cartId: %s", cartId)
	return &cart, nil
}

// Delete an item from cart by itemId
func (c *cDbService) DeleteItemFromCart(ctx context.Context, cartId string, itemId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing DeleteItemFromCart, cartId: %s, itemId: %s", cartId, itemId)

	filter := bson.M{"cartId": cartId}
	update := bson.M{"$pull": bson.M{"items": bson.M{"itemId": itemId}}}

	_, dbError := c.ucollection.UpdateOne(ctx, filter, update)
	if dbError != nil {
		logger.Error(dbError)
		return dbError
	}

	logger.Infof("Executed DeleteItemFromCart, cartId: %s, itemId: %s", cartId, itemId)
	return nil
}

// Delete all items from the cart
func (c *cDbService) DeleteAllItemsFromCart(ctx context.Context, cartId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing DeleteAllItemsFromCart, cartId: %s", cartId)

	filter := bson.M{"cartId": cartId}
	update := bson.M{"$set": bson.M{"items": []models.CartItem{}}}

	_, dbError := c.ucollection.UpdateOne(ctx, filter, update)
	if dbError != nil {
		logger.Error(dbError)
		return dbError
	}

	logger.Infof("Executed DeleteAllItemsFromCart, cartId: %s", cartId)
	return nil
}
