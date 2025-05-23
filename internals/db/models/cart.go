package dbmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartItem represents an item in the cart
type CartItem struct {
	ItemID   string  `json:"item_id" bson:"item_id"`
	Quantity int     `json:"quantity" bson:"quantity"`
	Price    float64 `json:"price" bson:"price"`
	Name     string  `json:"name" bson:"name"`
}

// Cart represents the structure of a user's cart
type Cart struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Items      []CartItem         `json:"items" bson:"items"`
	TotalPrice float64            `json:"totalprice"`
}
