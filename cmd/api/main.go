package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"erp-api/internal/delivery/http/client"
	"erp-api/internal/delivery/http/product"
	"erp-api/internal/delivery/http/quote"
	"erp-api/internal/delivery/http/reports"
	settingsHandler "erp-api/internal/delivery/http/settings"
	"erp-api/internal/delivery/http/tenant"
	"erp-api/internal/delivery/http/user"
	"erp-api/internal/infra/container"
	"erp-api/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	container := container.NewContainer()
	if err := container.Initialize(); err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	router := setupRouter(container)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(container *container.Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:3001", 
		"http://localhost:5173",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
		"http://127.0.0.1:5173",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true

	router.Use(cors.New(corsConfig))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authMiddleware := middleware.NewAuthMiddleware(container.GetJWTManager())

	api := router.Group("/api/v1")
	{
		tenants := api.Group("/tenants")
		{
			tenants.POST("", tenant.NewHandler(container.GetTenantUseCase()).Create)
			tenants.GET("/:id", tenant.NewHandler(container.GetTenantUseCase()).GetByID)
			tenants.PUT("/:id", tenant.NewHandler(container.GetTenantUseCase()).Update)
			tenants.DELETE("/:id", tenant.NewHandler(container.GetTenantUseCase()).Delete)
			tenants.GET("", tenant.NewHandler(container.GetTenantUseCase()).List)
			tenants.GET("/count", tenant.NewHandler(container.GetTenantUseCase()).Count)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/login", user.NewHandler(container.GetUserUseCase()).Login)
			auth.POST("/refresh", user.NewHandler(container.GetUserUseCase()).RefreshToken)
		}

		users := api.Group("/users")
		{
			users.POST("/register", user.NewHandler(container.GetUserUseCase()).Register)
			users.GET("/profile", authMiddleware.Authenticate(), user.NewHandler(container.GetUserUseCase()).GetProfile)
			users.GET("/count", authMiddleware.RequireRole("admin"), user.NewHandler(container.GetUserUseCase()).Count)
			users.GET("/:id", authMiddleware.RequireRole("admin"), user.NewHandler(container.GetUserUseCase()).GetByID)
			users.PUT("/:id", authMiddleware.RequireRole("admin"), user.NewHandler(container.GetUserUseCase()).Update)
			users.DELETE("/:id", authMiddleware.RequireRole("admin"), user.NewHandler(container.GetUserUseCase()).Delete)
			users.GET("", authMiddleware.RequireRole("admin"), user.NewHandler(container.GetUserUseCase()).List)
		}

		clients := api.Group("/clients")
		{
			clients.POST("", authMiddleware.Authenticate(), client.NewHandler(container.GetClientUseCase()).Create)
			clients.GET("/:id", authMiddleware.Authenticate(), client.NewHandler(container.GetClientUseCase()).GetByID)
			clients.PUT("/:id", authMiddleware.Authenticate(), client.NewHandler(container.GetClientUseCase()).Update)
			clients.DELETE("/:id", authMiddleware.Authenticate(), client.NewHandler(container.GetClientUseCase()).Delete)
			clients.GET("", authMiddleware.Authenticate(), client.NewHandler(container.GetClientUseCase()).List)
			clients.GET("/count", authMiddleware.Authenticate(), client.NewHandler(container.GetClientUseCase()).Count)
		}

		products := api.Group("/products")
		{
			products.POST("", authMiddleware.Authenticate(), product.NewHandler(container.GetProductUseCase()).Create)
			products.GET("/:id", authMiddleware.Authenticate(), product.NewHandler(container.GetProductUseCase()).GetByID)
			products.PUT("/:id", authMiddleware.Authenticate(), product.NewHandler(container.GetProductUseCase()).Update)
			products.DELETE("/:id", authMiddleware.Authenticate(), product.NewHandler(container.GetProductUseCase()).Delete)
			products.GET("", authMiddleware.Authenticate(), product.NewHandler(container.GetProductUseCase()).List)
			products.GET("/count", authMiddleware.Authenticate(), product.NewHandler(container.GetProductUseCase()).Count)
		}

		quotes := api.Group("/quotes")
		{
			quotes.POST("", authMiddleware.Authenticate(), quote.NewHandler(container.GetQuoteUseCase()).Create)
			quotes.GET("/:id", authMiddleware.Authenticate(), quote.NewHandler(container.GetQuoteUseCase()).GetByID)
			quotes.PUT("/:id", authMiddleware.Authenticate(), quote.NewHandler(container.GetQuoteUseCase()).Update)
			quotes.DELETE("/:id", authMiddleware.Authenticate(), quote.NewHandler(container.GetQuoteUseCase()).Delete)
			quotes.GET("", authMiddleware.Authenticate(), quote.NewHandler(container.GetQuoteUseCase()).List)
			quotes.GET("/count", authMiddleware.Authenticate(), quote.NewHandler(container.GetQuoteUseCase()).Count)
			quotes.PUT("/:id/status", authMiddleware.Authenticate(), quote.NewHandler(container.GetQuoteUseCase()).UpdateStatus)
		}

		settings := api.Group("/settings")
		{
			settings.GET("", authMiddleware.Authenticate(), settingsHandler.NewHandler(container.GetSettingsUseCase()).Get)
			settings.PUT("", authMiddleware.Authenticate(), settingsHandler.NewHandler(container.GetSettingsUseCase()).Update)
		}

		reportsGroup := api.Group("/reports")
		{
			reportsGroup.GET("/export", reports.NewHandler(container.GetProductUseCase()).Export)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})	

	router.GET("/health/readiness", func(c *gin.Context) {
		db, err := container.DB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"message": "Database connection failed",
			})
			return
		}
	
		if err := db.PingContext(context.Background()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"message": "Database connection failed", 
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return router
}