package services

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/db"
	dbmodel "Jevan/internals/db/models"
	"Jevan/internals/models"
	"context"
	"encoding/json"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *models.Product) (string, error)
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product, id string) error
	GetProductById(ctx context.Context, id string) (*models.Product, error)
	DeleteProductById(ctx context.Context, id string) error
}

type productService struct {
	db db.ProductDbService
}

func NewProductService(db db.ProductDbService) ProductService {
	return &productService{db: db}
}

func (p *productService) CreateProduct(ctx context.Context, product *models.Product) (string, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing CreateProduct product %v", product)
	var schema *dbmodel.ProductSchema
	b, _ := json.Marshal(product)
	_ = json.Unmarshal(b, &schema)
	return p.db.CreateProduct(ctx, schema)
}

func (p *productService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	return p.db.GetAllProducts(ctx)
}

func (p *productService) UpdateProduct(ctx context.Context, product *models.Product, id string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing UpdateProduct id: %s", id)
	var schema *dbmodel.ProductSchema
	b, _ := json.Marshal(product)
	_ = json.Unmarshal(b, &schema)
	return p.db.UpdateProduct(ctx, schema, id)
}

func (p *productService) GetProductById(ctx context.Context, id string) (*models.Product, error) {
	return p.db.GetProductById(ctx, id)
}

func (p *productService) DeleteProductById(ctx context.Context, id string) error {
	return p.db.DeleteProductById(ctx, id)
}
