package main

import (
	"Jevan/apis"
	"Jevan/apis/middlewares"
	"Jevan/internals/db"
	"Jevan/internals/services"

	_ "Jevan/apis/docs"
	"Jevan/commons/apploggers"
	"Jevan/configs"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Jevan - Mess Management API
// @version 1.0
// @description Backend APIs for Jevan mess application using Echo.
// @contact.name API Support
// @contact.email support@jevan.app
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	userDbService := db.NewUserDbService(configs.AppConfig.DbClient)

	// Initialize services
	productService := services.NewProductService(productDbService)
	cartService := services.NewCartService(cartDbService)
	orderService := services.NewOrderService(orderDbService)
	userService := services.NewUserService(userDbService)

	// Controllers
	productController := apis.NewProductController(productService)
	cartController := apis.NewCartController(cartService)
	orderController := apis.NewOrderController(orderService)
	userController := apis.NewUserController(userService)
	authController := apis.NewAuthController(userService)

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	jwtMiddleware := middlewares.JWTMiddleware()

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Jevan API is healthy!")
	})

	// Auth
	e.POST("/login", authController.Login)
	e.POST("/register", authController.Register)

	// Admin-only endpoints
	admin := e.Group("/admin")
	admin.Use(jwtMiddleware, middlewares.AdminOnly)
	admin.PUT("/users/:id/role", authController.UpdateUserRole)

	// Public Routes
	e.GET("/users", userController.GetUsers)
	e.GET("/users/:id", userController.GetUserById)

	// Auth-Protected User Actions
	userPrivate := e.Group("/users", jwtMiddleware)
	userPrivate.DELETE(":id", userController.DeleteUserById)
	userPrivate.POST("", userController.CreateUser)
	userPrivate.PATCH(":id", userController.UpdateUser)

	// Product Routes
	productPublic := e.Group("/products")
	productPublic.GET("", productController.GetAllProducts)
	productPublic.GET(":id", productController.GetProductById)

	productPrivate := e.Group("/products", jwtMiddleware)
	productPrivate.POST("", productController.CreateProduct)
	productPrivate.PUT(":id", productController.UpdateProduct)
	productPrivate.DELETE(":id", productController.DeleteProductById)

	// Cart Routes
	cart := e.Group("/cart", jwtMiddleware)
	cart.POST("", cartController.UpdateCart)
	cart.GET(":id", cartController.GetCartItemsById)
	cart.DELETE(":id/all", cartController.DeleteAllItems)
	cart.PUT(":cartId/item/:itemId", cartController.UpdateItemQuantity)

	// Order Routes
	order := e.Group("/orders", jwtMiddleware)
	order.POST("", orderController.CreateOrder)
	order.GET("", orderController.GetAllOrders)
	order.GET(":id", orderController.GetOrderById)
	order.PUT(":id", orderController.UpdateOrder)

	logger.Infof("Starting Jevan API server on port %s", configs.AppConfig.HttpPort)
	e.Logger.Fatal(e.Start(":" + configs.AppConfig.HttpPort))
}
