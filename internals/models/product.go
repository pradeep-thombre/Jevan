package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Category    string             `json:"category" bson:"category"`
	ImageURL    string             `json:"image" bson:"image"`
	IsAvailable bool               `json:"isAvailable" bson:"isAvailable"`
	Rating      float64            `json:"rating" bson:"rating"`
	Type        string             `json:"type" bson:"type"`
	MealTime    string             `json:"mealTime" bson:"mealTime"`
}
