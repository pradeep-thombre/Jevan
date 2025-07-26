package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `bson:"firstName" json:"firstName" validate:"required"`
	LastName  string             `bson:"lastName" json:"lastName" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	CartId    string             `json:"cart_id"`
	Type      string             `json:"type"`
	Age       int                `json:"age"`
	IsActive  bool               `json:"is_active"`
}

type UserDetails struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstName" json:"firstName" validate:"required"`
	LastName  string             `bson:"lastName" json:"lastName" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password,omitempty" validate:"required,min=6"`
	Role      string             `bson:"role" json:"role" validate:"omitempty,oneof=admin user"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin user"`
}
