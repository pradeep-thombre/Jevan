package apis

import (
	"Jevan/commons/apploggers"
	"Jevan/internals/models"
	"Jevan/internals/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// @Summary Create Product
// @Description Creates a new product
// @Tags Product
// @Accept json
// @Produce json
// @Param product body models.Product true "Product Info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /products [post]
func (pc *ProductController) CreateProduct(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	logger.Info("Received request to create product")

	var product models.Product
	if err := c.Bind(&product); err != nil {
		logger.Error("Failed to bind product request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	id, err := pc.productService.CreateProduct(c.Request().Context(), &product)
	if err != nil {
		logger.Error("Failed to create product: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create product"})
	}

	logger.Infof("Product created with ID: %s", id)
	return c.JSON(http.StatusCreated, map[string]string{"productId": id})
}

// @Summary Get All Products
// @Description Retrieves all products
// @Tags Product
// @Produce json
// @Success 200 {array} models.Product
// @Router /products [get]
func (pc *ProductController) GetAllProducts(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	logger.Info("Received request to get all products")

	products, err := pc.productService.GetAllProducts(c.Request().Context())
	if err != nil {
		logger.Error("Failed to get all products: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products"})
	}

	logger.Infof("Fetched %d products", len(products))
	return c.JSON(http.StatusOK, products)
}

// @Summary Update Product
// @Description Updates an existing product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product Info"
// @Success 200 {object} map[string]string
// @Router /products/{id} [put]
func (pc *ProductController) UpdateProduct(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	id := c.Param("id")
	logger.Infof("Received request to update product with ID: %s", id)

	var product models.Product
	if err := c.Bind(&product); err != nil {
		logger.Error("Failed to bind product update request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := pc.productService.UpdateProduct(c.Request().Context(), &product, id)
	if err != nil {
		logger.Error("Failed to update product: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	}

	logger.Infof("Successfully updated product with ID: %s", id)
	return c.JSON(http.StatusOK, map[string]string{"message": "Product updated successfully"})
}

// @Summary Get Product by ID
// @Description Retrieves a product by its ID
// @Tags Product
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func (pc *ProductController) GetProductById(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	id := c.Param("id")
	logger.Infof("Received request to get product by ID: %s", id)

	product, err := pc.productService.GetProductById(c.Request().Context(), id)
	if err != nil {
		logger.Error("Failed to get product by ID: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch product"})
	}

	logger.Infof("Fetched product with ID: %s", id)
	return c.JSON(http.StatusOK, product)
}

// @Summary Delete Product by ID
// @Description Deletes a product by its ID
// @Tags Product
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Router /products/{id} [delete]
func (pc *ProductController) DeleteProductById(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	id := c.Param("id")
	logger.Infof("Received request to delete product with ID: %s", id)

	err := pc.productService.DeleteProductById(c.Request().Context(), id)
	if err != nil {
		logger.Error("Failed to delete product: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete product"})
	}

	logger.Infof("Successfully deleted product with ID: %s", id)
	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
