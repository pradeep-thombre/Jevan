package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     string             `json:"userId" validate:"required"`
	Items      []OrderItem        `json:"items" validate:"required"`
	TotalPrice float64            `json:"totalPrice"`
	Status     string             `json:"status"`    // e.g., "Order Placed", "Ready", "Preparing", "Delivered", "Shipped"
	OrderedAt  int64              `json:"orderedAt"` // Unix timestamp
	UpdatedAt  int64              `json:"updatedAt"`
}

type OrderItem struct {
	ItemID   string `json:"itemId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}
