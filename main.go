package main

import (
	"Jevan/apis"
	"Jevan/internals/db"
	"Jevan/internals/services"

	_ "Jevan/apis/docs"
	"Jevan/commons/apploggers"
	"Jevan/configs"
	"context"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Jevan - Mess Management API
// @version 1.0
// @description Backend APIs for Jevan mess application using Echo.
// @contact.name API Support
// @contact.email support@jevan.app
// @host localhost:3000
// @BasePath /
func main() {
	// Set up logging and context
	ctx, logger := apploggers.NewLoggerWithCorrelationid(context.Background(), "")

	// Load configuration
	if err := configs.NewApplicationConfig(ctx); err != nil {
		logger.Errorf("Failed to load app config: %v", err)
	}

	// Initialize DB
	cartDbService := db.NewCartDbService(configs.AppConfig.DbClient)
	orderDbService := db.NewOrderDbService(configs.AppConfig.DbClient)
	productDbService := db.NewProductDbService(configs.AppConfig.DbClient)
	dbservice := db.NewUserDbService(configs.AppConfig.DbClient)

	// Initialize services
	productService := services.NewProductService(productDbService)
	cartService := services.NewCartService(cartDbService)
	orderService := services.NewOrderService(orderDbService)
	eventService := services.NewUserEventService(dbservice)

	// Echo instance
	e := echo.New()

	// user api Routes
	userController := apis.NewUserController(eventService)
	e.GET("/users", userController.GetUsers)
	e.GET("/users/:id", userController.GetUserById)
	e.DELETE("/users/:id", userController.DeleteUserById)
	e.POST("/users", userController.CreateUser)
	e.PATCH("/users/:id", userController.UpdateUser)

	// Product routes
	productController := apis.NewProductController(productService)
	e.POST("/products", productController.CreateProduct)
	e.GET("/products", productController.GetAllProducts)
	e.PUT("/products/:id", productController.UpdateProduct)
	e.GET("/products/:id", productController.GetProductById)
	e.DELETE("/products/:id", productController.DeleteProductById)

	// Cart routes
	cartController := apis.NewCartController(cartService)
	e.POST("/cart", cartController.UpdateCart)
	e.GET("/cart/:id", cartController.GetCartItemsById)
	e.DELETE("/cart/:id/all", cartController.DeleteAllItems)
	e.PUT("/cart/:cartId/item/:itemId", cartController.UpdateItemQuantity)

	// Order routes
	orderController := apis.NewOrderController(orderService)
	e.POST("/orders", orderController.CreateOrder)
	e.GET("/orders/:id", orderController.GetOrderById)
	e.PUT("/orders/:id", orderController.UpdateOrder)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Jevan API is healthy!")
	})

	// Start server
	logger.Infof("Starting Jevan API server on port %s", configs.AppConfig.HttpPort)
	e.Logger.Fatal(e.Start(":" + configs.AppConfig.HttpPort))
}
