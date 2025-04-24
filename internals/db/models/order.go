package dbmodels

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderSchema struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Items      []OrderItemSchema  `bson:"items"`
	TotalPrice float64            `bson:"totalprice"`
	Status     string             `bson:"status"`     // e.g., "pending", "confirmed", "delivered"
	OrderedAt  int64              `bson:"ordered_at"` // Unix timestamp
}

type OrderItemSchema struct {
	ProductID primitive.ObjectID `bson:"product_id"`
	Name      string             `bson:"name"`
	Price     float64            `bson:"price"`
	Quantity  int                `bson:"quantity"`
}
