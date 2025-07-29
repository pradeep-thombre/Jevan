package services

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/db"
	"Jevan/internals/models"
	"context"
)

type CartService interface {
	UpdateCart(ctx context.Context, cart *models.Cart) error
	GetCartItemsById(ctx context.Context, cartId string) (*models.Cart, error)
	DeleteAllItems(ctx context.Context, cartId string) error
}

type cartService struct {
	dbservice db.CartDbService
}

func NewCartService(dbservice db.CartDbService) CartService {
	return &cartService{dbservice: dbservice}
}

func (c *cartService) GetCartItemsById(ctx context.Context, cartId string) (*models.Cart, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing GetCartItemsById, cartId: %s", cartId)

	cart, err := c.dbservice.GetCartById(ctx, cartId)
	if err != nil {
		logger.Errorf("Failed to get cart for id %s: %v", cartId, err)
		return nil, err
	}

	logger.Infof("Fetched cart successfully for id: %s", cartId)
	return cart, nil
}

func (c *cartService) DeleteAllItems(ctx context.Context, cartId string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing DeleteAllItems, cartId: %s", cartId)

	err := c.dbservice.DeleteAllItemsFromCart(ctx, cartId)
	if err != nil {
		logger.Errorf("Failed to delete items for cart %s: %v", cartId, err)
		return err
	}

	logger.Infof("Deleted all items from cart %s successfully", cartId)
	return nil
}

func (cs *cartService) UpdateCart(ctx context.Context, cart *models.Cart) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing UpdateCart, cartId: %s", cart.ID)

	err := cs.dbservice.SaveCart(ctx, cart)
	if err != nil {
		logger.Errorf("Failed to update cart %s: %v", cart.ID, err)
		return err
	}

	logger.Infof("Cart %s updated successfully", cart.ID)
	return nil
}
