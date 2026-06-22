// Package main PriceCompare API
// @title PriceCompare API
// @version 1.0
// @description API for comparing product prices across stores
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	"github.com/joho/godotenv"
	"pricecompare/internal/config"
	"pricecompare/internal/database"
	"pricecompare/internal/handlers"
	"pricecompare/internal/middleware"
	"pricecompare/internal/repositories"
	"pricecompare/internal/seed"
	"pricecompare/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "pricecompare/docs"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if err := seed.Run(db); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}

	if err := seed.UpdateImages(db); err != nil {
		log.Fatalf("Image update failed: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	storeRepo := repositories.NewStoreRepository(db)
	productRepo := repositories.NewProductRepository(db)
	priceRepo := repositories.NewPriceRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	authService := services.NewAuthService(userRepo, cfg)
	productService := services.NewProductService(productRepo, priceRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	storeService := services.NewStoreService(storeRepo)
	priceService := services.NewPriceService(priceRepo)
	cartService := services.NewCartService(cartRepo, priceRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, priceRepo)

	h := handlers.NewHandlers(
		authService, productService, categoryService, storeService,
		priceService, cartService, orderService,
	)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.GET("/home", h.GetHome)

		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)

		api.GET("/categories", h.GetCategories)
		api.GET("/categories/:id", h.GetCategory)

		api.GET("/stores", h.GetStores)
		api.GET("/stores/:id", h.GetStore)

		api.GET("/products", h.GetProducts)
		api.GET("/products/:id", h.GetProduct)
		api.GET("/prices", h.GetPrices)

		cartGroup := api.Group("/cart")
		cartGroup.Use(middleware.OptionalAuthMiddleware(authService))
		{
			cartGroup.GET("", h.GetCart)
			cartGroup.POST("", h.AddToCart)
			cartGroup.PUT("/:id", h.UpdateCartItem)
			cartGroup.DELETE("/:id", h.RemoveCartItem)
			cartGroup.DELETE("", h.ClearCart)
		}

		authGroup := api.Group("")
		authGroup.Use(middleware.AuthMiddleware(authService))
		{
			authGroup.GET("/users/me", h.GetProfile)
			authGroup.POST("/orders", h.CreateOrder)
			authGroup.GET("/orders/my", h.GetMyOrders)
		}

		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware(authService), middleware.AdminMiddleware())
		{
			adminGroup.POST("/categories", h.CreateCategory)
			adminGroup.PUT("/categories/:id", h.UpdateCategory)
			adminGroup.DELETE("/categories/:id", h.DeleteCategory)

			adminGroup.POST("/stores", h.CreateStore)
			adminGroup.PUT("/stores/:id", h.UpdateStore)
			adminGroup.DELETE("/stores/:id", h.DeleteStore)

			adminGroup.POST("/products", h.CreateProduct)
			adminGroup.PUT("/products/:id", h.UpdateProduct)
			adminGroup.DELETE("/products/:id", h.DeleteProduct)

			adminGroup.POST("/prices", h.CreatePrice)
			adminGroup.PUT("/prices/:id", h.UpdatePrice)
			adminGroup.DELETE("/prices/:id", h.DeletePrice)

			adminGroup.GET("/orders", h.GetAllOrders)
			adminGroup.PATCH("/orders/:id/status", h.UpdateOrderStatus)
		}
	}

	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
