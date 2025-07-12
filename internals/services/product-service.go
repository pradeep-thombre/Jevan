package services

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/db"
	"Jevan/internals/models"
	"context"
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
	logger.Infof("Executing CreateProduct: %v", product)

	productId, err := p.db.CreateProduct(ctx, product)
	if err != nil {
		logger.Errorf("Failed to create product: %v", err)
		return "", err
	}

	logger.Infof("Product created successfully: %s", productId)
	return productId, nil
}

func (p *productService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Info("Executing GetAllProducts")

	products, err := p.db.GetAllProducts(ctx)
	if err != nil {
		logger.Errorf("Failed to fetch products: %v", err)
		return nil, err
	}

	logger.Infof("Fetched %d products", len(products))
	return products, nil
}

func (p *productService) UpdateProduct(ctx context.Context, product *models.Product, id string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing UpdateProduct id: %s, payload: %v", id, product)

	err := p.db.UpdateProduct(ctx, product, id)
	if err != nil {
		logger.Errorf("Failed to update product %s: %v", id, err)
		return err
	}

	logger.Infof("Product %s updated successfully", id)
	return nil
}

func (p *productService) GetProductById(ctx context.Context, id string) (*models.Product, error) {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing GetProductById for id: %s", id)

	product, err := p.db.GetProductById(ctx, id)
	if err != nil {
		logger.Errorf("Failed to fetch product %s: %v", id, err)
		return nil, err
	}

	logger.Infof("Fetched product %s successfully", id)
	return product, nil
}

func (p *productService) DeleteProductById(ctx context.Context, id string) error {
	logger := apploggers.GetLoggerWithCorrelationid(ctx)
	logger.Infof("Executing DeleteProductById for id: %s", id)

	err := p.db.DeleteProductById(ctx, id)
	if err != nil {
		logger.Errorf("Failed to delete product %s: %v", id, err)
		return err
	}

	logger.Infof("Product %s deleted successfully", id)
	return nil
}
