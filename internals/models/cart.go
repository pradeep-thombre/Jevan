package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// CartItem represents an item in the cart
type CartItem struct {
	ItemID   string  `json:"item_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Name     string  `json:"name"`
}

// Cart represents the structure of a user's cart
type Cart struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Items      []CartItem         `json:"items"`
	TotalPrice float64            `json:"totalprice"`
}
