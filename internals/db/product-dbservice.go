package db

import (
	"Jevan/commons/appdb"
	"Jevan/commons/apploggers"
	"Jevan/configs"
	dbmodel "Jevan/internals/db/models"
	"Jevan/internals/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductDbService interface {
	CreateProduct(ctx context.Context, product *dbmodel.ProductSchema) (string, error)
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, product *dbmodel.ProductSchema, id string) error
	GetProductById(ctx context.Context, id string) (*models.Product, error)
	DeleteProductById(ctx context.Context, id string) error
}

type productDb struct {
	collection appdb.DatabaseCollection
}

func NewProductDbService(client appdb.DatabaseClient) ProductDbService {
	return &productDb{
		collection: client.Collection(configs.MONGO_PRODUCTS_COLLECTION),
	}
}

func (p *productDb) CreateProduct(ctx context.Context, product *dbmodel.ProductSchema) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Creating product: %+v", product)

	result, err := p.collection.InsertOne(ctx, product)
	if err != nil {
		logger.Error("Failed to insert product: ", err)
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	logger.Infof("Product created with ID: %s", id)
	return id, nil
}

func (p *productDb) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Info("Fetching all products")

	var products []*models.Product
	err := p.collection.Find(ctx, bson.M{}, &options.FindOptions{}, &products)
	if err != nil {
		logger.Error("Failed to fetch products: ", err)
		return nil, err
	}

	logger.Infof("Fetched %d products", len(products))
	return products, nil
}

func (p *productDb) UpdateProduct(ctx context.Context, product *dbmodel.ProductSchema, id string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Updating product with ID: %s", id)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Errorf("Invalid product ID: %s", id)
		return fmt.Errorf("invalid id: %s", id)
	}

	_, err = p.collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": product})
	if err != nil {
		logger.Error("Failed to update product: ", err)
		return err
	}

	logger.Infof("Successfully updated product with ID: %s", id)
	return nil
}

func (p *productDb) GetProductById(ctx context.Context, id string) (*models.Product, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Fetching product by ID: %s", id)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Errorf("Invalid product ID format: %s", id)
		return nil, fmt.Errorf("invalid id: %s", id)
	}

	var product *models.Product
	err = p.collection.FindOne(ctx, bson.M{"_id": objId}, &product)
	if err != nil {
		logger.Error("Failed to fetch product: ", err)
		return nil, err
	}

	logger.Infof("Fetched product: %+v", product)
	return product, nil
}

func (p *productDb) DeleteProductById(ctx context.Context, id string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Deleting product with ID: %s", id)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Errorf("Invalid product ID: %s", id)
		return fmt.Errorf("invalid id: %s", id)
	}

	_, err = p.collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		logger.Error("Failed to delete product: ", err)
		return err
	}

	logger.Infof("Successfully deleted product with ID: %s", id)
	return nil
}
