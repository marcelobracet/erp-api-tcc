package container

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	userDomain "erp-api/internal/domain/user"
	"erp-api/internal/infra/repository"
	userUseCase "erp-api/internal/usecase/user"
	"erp-api/pkg/auth"

	_ "github.com/lib/pq"
)

type Container struct {
	DB          *sql.DB
	UserRepo    userDomain.Repository
	UserUseCase userUseCase.UseCaseInterface
	JWTManager  *auth.JWTManager
	PassHasher  *auth.PasswordHasher
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Initialize() error {
	if err := c.initializeDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := c.initializeAuth(); err != nil {
		return fmt.Errorf("failed to initialize auth: %w", err)
	}

	if err := c.initializeRepositories(); err != nil {
		return fmt.Errorf("failed to initialize repositories: %w", err)
	}

	if err := c.initializeUseCases(); err != nil {
		return fmt.Errorf("failed to initialize use cases: %w", err)
	}

	return nil
}

func (c *Container) initializeDatabase() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "erp_user"
	}
	if dbPassword == "" {
		dbPassword = "erp_password"
	}
	if dbName == "" {
		dbName = "erp_db"
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	c.DB = db
	log.Println("Database initialized successfully")
	return nil
}

func (c *Container) initializeRepositories() error {
	if c.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	c.UserRepo = repository.NewUserRepository(c.DB)
	log.Println("Repositories initialized successfully")
	return nil
}

func (c *Container) initializeUseCases() error {
	if c.UserRepo == nil {
		return fmt.Errorf("repositories not initialized")
	}

	c.UserUseCase = userUseCase.NewUseCase(c.UserRepo, c.JWTManager, c.PassHasher)
	log.Printf("Use cases initialized successfully - UserRepo: %v, JWTManager: %v, PassHasher: %v", c.UserRepo != nil, c.JWTManager != nil, c.PassHasher != nil)
	return nil
}

func (c *Container) initializeAuth() error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-here"
	}

	accessExpiryStr := os.Getenv("JWT_EXPIRATION")
	if accessExpiryStr == "" {
		accessExpiryStr = "24h"
	}

	refreshExpiryStr := os.Getenv("JWT_REFRESH_EXPIRATION")
	if refreshExpiryStr == "" {
		refreshExpiryStr = "168h"
	}

	accessExpiry, err := time.ParseDuration(accessExpiryStr)
	if err != nil {
		accessExpiry = 24 * time.Hour
	}

	refreshExpiry, err := time.ParseDuration(refreshExpiryStr)
	if err != nil {
		refreshExpiry = 168 * time.Hour
	}

	c.JWTManager = auth.NewJWTManager(jwtSecret, accessExpiry, refreshExpiry)
	c.PassHasher = auth.DefaultPasswordHasher()

	log.Printf("Auth initialized successfully - JWTManager: %v, PassHasher: %v", c.JWTManager != nil, c.PassHasher != nil)
	return nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

func (c *Container) GetUserRepository() userDomain.Repository {
	return c.UserRepo
}

func (c *Container) GetUserUseCase() userUseCase.UseCaseInterface {
	return c.UserUseCase
}

func (c *Container) GetJWTManager() *auth.JWTManager {
	return c.JWTManager
}

func (c *Container) GetPassHasher() *auth.PasswordHasher {
	return c.PassHasher
} 