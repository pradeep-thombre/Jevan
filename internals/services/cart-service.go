package services

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/db"
	"Jevan/internals/models"
	"context"
	"fmt"
)

type CartService interface {
	UpdateCart(ctx context.Context, cart *models.Cart) error
	UpdateItemQuantity(ctx context.Context, cartId string, itemId string, quantity int) (*models.Cart, error)
	GetCartItemsById(ctx context.Context, cartId string) (*models.Cart, error)
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

func (cs *cartService) UpdateCart(ctx context.Context, cart *models.Cart) error {
	return cs.dbservice.SaveCart(ctx, cart)
}

func (cs *cartService) UpdateItemQuantity(ctx context.Context, cartId, itemId string, quantity int) (*models.Cart, error) {
	cart, err := cs.GetCartItemsById(ctx, cartId)
	if err != nil {
		return nil, err
	}

	var updatedItems []models.CartItem
	var total float64

	for _, item := range cart.Items {
		if item.ItemID == itemId {
			if quantity == 0 {
				continue
			}
			item.Quantity = quantity
		}
		total += float64(item.Quantity) * item.Price
		updatedItems = append(updatedItems, item)
	}

	cart.Items = updatedItems
	cart.TotalPrice = total
	err = cs.dbservice.SaveCart(ctx, cart)
	return cart, err
}
