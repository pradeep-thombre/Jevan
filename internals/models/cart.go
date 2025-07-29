package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// CartItem represents an item in the cart
type CartItem struct {
	ItemID   string `json:"itemId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

// Cart represents the structure of a user's cart
type Cart struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Items      []CartItem         `json:"items" validate:"required"`
	TotalPrice float64            `json:"totalPrice"`
}
