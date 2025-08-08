package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"erp-api/internal/delivery/http/user"
	"erp-api/internal/infra/container"
	"erp-api/pkg/middleware"

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

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authMiddleware := middleware.NewAuthMiddleware(container.GetJWTManager())

	api := router.Group("/api/v1")
	{
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
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return router
} 