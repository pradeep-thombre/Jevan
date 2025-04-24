package dbmodels

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserSchema struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	CartId   string             `json:"cart_id"`
	Type     string             `json:"type"`
	Age      int                `json:"age"`
	IsActive bool               `json:"is_active"`
}
