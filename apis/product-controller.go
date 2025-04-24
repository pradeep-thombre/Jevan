package apis

import (
	"Jevan/commons"
	"Jevan/commons/apploggers"
	"Jevan/internals/models"
	"Jevan/internals/services"
	"net/http"
	"strings"

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
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /products [post]
func (pc *ProductController) CreateProduct(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	logger.Info("Received request to create product")

	var product models.Product
	if err := c.Bind(&product); err != nil {
		logger.Error("Invalid request body: ", err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Invalid request body", nil))
	}

	id, err := pc.productService.CreateProduct(c.Request().Context(), &product)
	if err != nil {
		logger.Error("Failed to create product: ", err)
		return c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to create product", nil))
	}

	logger.Infof("Product created with ID: %s", id)
	return c.JSON(http.StatusCreated, map[string]string{"productId": id})
}

// @Summary Get All Products
// @Description Retrieves all products
// @Tags Product
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func (pc *ProductController) GetAllProducts(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	logger.Info("Received request to get all products")

	products, err := pc.productService.GetAllProducts(c.Request().Context())
	if err != nil {
		logger.Error("Failed to fetch products: ", err)
		return c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to fetch products", nil))
	}

	logger.Infof("Fetched %d products", len(products))
	return c.JSON(http.StatusOK, map[string]interface{}{
		"total":    len(products),
		"products": products,
	})
}

// @Summary Update Product
// @Description Updates an existing product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product Info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /products/{id} [put]
func (pc *ProductController) UpdateProduct(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	id := c.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}

	logger.Infof("Received request to update product with ID: %s", id)

	var product models.Product
	if err := c.Bind(&product); err != nil {
		logger.Error("Invalid request body: ", err)
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("Invalid request body", nil))
	}

	if err := pc.productService.UpdateProduct(c.Request().Context(), &product, id); err != nil {
		logger.Error("Failed to update product: ", err)
		return c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to update product", nil))
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
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /products/{id} [get]
func (pc *ProductController) GetProductById(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	id := c.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}

	logger.Infof("Received request to get product by ID: %s", id)

	product, err := pc.productService.GetProductById(c.Request().Context(), id)
	if err != nil {
		logger.Error("Failed to fetch product: ", err)
		return c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to fetch product", nil))
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
// @Failure 400 {object} commons.ApiErrorResponsePayload
// @Router /products/{id} [delete]
func (pc *ProductController) DeleteProductById(c echo.Context) error {
	logger := apploggers.GetLoggerWithCorrelationid(c.Request().Context())
	id := c.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		logger.Error("'id' is required")
		return c.JSON(http.StatusBadRequest, commons.ApiErrorResponse("'id' is required", nil))
	}

	logger.Infof("Received request to delete product with ID: %s", id)

	if err := pc.productService.DeleteProductById(c.Request().Context(), id); err != nil {
		logger.Error("Failed to delete product: ", err)
		return c.JSON(http.StatusInternalServerError, commons.ApiErrorResponse("Failed to delete product", nil))
	}

	logger.Infof("Successfully deleted product with ID: %s", id)
	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
