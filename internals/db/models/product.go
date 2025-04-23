package dbmodels

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductSchema struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	Category    string             `bson:"category"`
	ImageURL    string             `bson:"image_url"`
	IsAvailable bool               `bson:"is_available"`
}
