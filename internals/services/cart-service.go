package services

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/db"
	"Jevan/internals/models"
	"context"
	"fmt"
)

type CartService interface {
	AddItemToCart(ctx context.Context, cartId string, item *models.CartItem) error
	GetCartItemsById(ctx context.Context, cartId string) (*models.Cart, error)
	DeleteItemsFromCart(ctx context.Context, cartId string, itemId string) error
	DeleteAllItems(ctx context.Context, cartId string) error
}

type cartService struct {
	dbservice db.CartDbService
}

func NewCartService(dbservice db.CartDbService) CartService {
	return &cartService{
		dbservice: dbservice,
	}
}

func (c *cartService) AddItemToCart(ctx context.Context, cartId string, item *models.CartItem) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing AddItemToCart, cartId: %s", cartId)

	err := c.dbservice.AddToCart(ctx, cartId, item)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to add item to cart: %s", err))
		return err
	}

	logger.Infof("Executed AddItemToCart, cartId: %s", cartId)
	return nil
}

func (c *cartService) GetCartItemsById(ctx context.Context, cartId string) (*models.Cart, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing GetCartItemsById, cartId: %s", cartId)

	cart, err := c.dbservice.GetCartById(ctx, cartId)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get cart: %s", err))
		return nil, err
	}

	logger.Infof("Executed GetCartItemsById, cartId: %s", cartId)
	return cart, nil
}

func (c *cartService) DeleteItemsFromCart(ctx context.Context, cartId string, itemId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing DeleteItemsFromCart, cartId: %s, itemId: %s", cartId, itemId)

	err := c.dbservice.DeleteItemFromCart(ctx, cartId, itemId)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to delete item from cart: %s", err))
		return err
	}

	logger.Infof("Executed DeleteItemsFromCart, cartId: %s, itemId: %s", cartId, itemId)
	return nil
}

func (c *cartService) DeleteAllItems(ctx context.Context, cartId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing DeleteAllItems, cartId: %s", cartId)

	err := c.dbservice.DeleteAllItemsFromCart(ctx, cartId)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to delete all items from cart: %s", err))
		return err
	}

	logger.Infof("Executed DeleteAllItems, cartId: %s", cartId)
	return nil
}
