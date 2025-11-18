package container

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	clientDomain "erp-api/internal/domain/client"
	productDomain "erp-api/internal/domain/product"
	quoteDomain "erp-api/internal/domain/quote"
	settingsDomain "erp-api/internal/domain/settings"
	tenantDomain "erp-api/internal/domain/tenant"
	userDomain "erp-api/internal/domain/user"
	"erp-api/internal/infra/repository"
	clientUseCase "erp-api/internal/usecase/client"
	productUseCase "erp-api/internal/usecase/product"
	quoteUseCase "erp-api/internal/usecase/quote"
	settingsUseCase "erp-api/internal/usecase/settings"
	tenantUseCase "erp-api/internal/usecase/tenant"
	userUseCase "erp-api/internal/usecase/user"
	"erp-api/pkg/auth"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Container struct {
	DB             *gorm.DB
	TenantRepo     tenantDomain.Repository
	TenantUseCase  tenantUseCase.UseCaseInterface
	UserRepo       userDomain.Repository
	UserUseCase    userUseCase.UseCaseInterface
	ClientRepo     clientDomain.Repository
	ClientUseCase  clientUseCase.UseCaseInterface
	ProductRepo    productDomain.Repository
	ProductUseCase productUseCase.UseCaseInterface
	QuoteRepo      quoteDomain.Repository
	QuoteItemRepo  quoteDomain.ItemRepository
	QuoteUseCase   quoteUseCase.UseCaseInterface
	SettingsRepo   settingsDomain.Repository
	SettingsUseCase settingsUseCase.UseCaseInterface
	JWTManager     *auth.JWTManager
	PassHasher     *auth.PasswordHasher
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Initialize() error {
	if err := c.initializeDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := c.runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	c.DB = db
	log.Println("Database initialized successfully")
	return nil
}

func (c *Container) runMigrations() error {
	if c.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Running database migrations...")

	// Habilitar extensão UUID se necessário
	if err := c.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: Could not create uuid-ossp extension: %v", err)
		// Tentar gen_random_uuid() que é nativo do PostgreSQL 13+
		if err := c.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
			log.Printf("Warning: Could not create pgcrypto extension: %v", err)
		}
	}

	// Executar AutoMigrate para todas as entidades na ordem correta
	// Ordem importa devido às foreign keys
	err := c.DB.AutoMigrate(
		&tenantDomain.Tenant{},
		&userDomain.User{},
		&clientDomain.Client{},
		&productDomain.Product{},
		&quoteDomain.Quote{},
		&quoteDomain.QuoteItem{},
		&settingsDomain.Settings{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	c.createForeignKeys()

	c.createGeneratedColumnsAndIndexes()

	log.Println("Database migrations completed successfully")
	return nil
}

func (c *Container) createForeignKeys() {
	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_users_tenant'
			) THEN
				ALTER TABLE users ADD CONSTRAINT fk_users_tenant 
				FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`)

	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_clients_tenant'
			) THEN
				ALTER TABLE clients ADD CONSTRAINT fk_clients_tenant 
				FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`)

	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_products_tenant'
			) THEN
				ALTER TABLE products ADD CONSTRAINT fk_products_tenant 
				FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`)

	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_quotes_tenant'
			) THEN
				ALTER TABLE quotes ADD CONSTRAINT fk_quotes_tenant 
				FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
			END IF;
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_quotes_client'
			) THEN
				ALTER TABLE quotes ADD CONSTRAINT fk_quotes_client 
				FOREIGN KEY (client_id) REFERENCES clients(id);
			END IF;
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_quotes_user'
			) THEN
				ALTER TABLE quotes ADD CONSTRAINT fk_quotes_user 
				FOREIGN KEY (user_id) REFERENCES users(id);
			END IF;
		END $$;
	`)

	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_quote_items_tenant'
			) THEN
				ALTER TABLE quote_items ADD CONSTRAINT fk_quote_items_tenant 
				FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
			END IF;
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_quote_items_quote'
			) THEN
				ALTER TABLE quote_items ADD CONSTRAINT fk_quote_items_quote 
				FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE CASCADE;
			END IF;
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_quote_items_product'
			) THEN
				ALTER TABLE quote_items ADD CONSTRAINT fk_quote_items_product 
				FOREIGN KEY (product_id) REFERENCES products(id);
			END IF;
		END $$;
	`)

	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_settings_tenant'
			) THEN
				ALTER TABLE settings ADD CONSTRAINT fk_settings_tenant 
				FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`)
}

func (c *Container) createGeneratedColumnsAndIndexes() {
	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'quote_items' AND column_name = 'total'
			) THEN
				ALTER TABLE quote_items ADD COLUMN total NUMERIC(12,2) 
				GENERATED ALWAYS AS (quantity * price) STORED;
			END IF;
		END $$;
	`)

	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_indexes WHERE indexname = 'idx_clients_document_tenant'
			) THEN
				CREATE UNIQUE INDEX idx_clients_document_tenant ON clients(document, tenant_id);
			END IF;
		END $$;
	`)

	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_indexes WHERE indexname = 'idx_settings_key_tenant'
			) THEN
				CREATE UNIQUE INDEX idx_settings_key_tenant ON settings(key, tenant_id);
			END IF;
		END $$;
	`)

	c.DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'clients_document_type_check'
			) THEN
				ALTER TABLE clients ADD CONSTRAINT clients_document_type_check 
				CHECK (document_type IN ('CPF', 'CNPJ'));
			END IF;
		END $$;
	`)
}

func (c *Container) initializeRepositories() error {
	if c.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	c.TenantRepo = repository.NewTenantRepository(c.DB)
	c.UserRepo = repository.NewUserRepository(c.DB)
	c.ClientRepo = repository.NewClientRepository(c.DB)
	c.ProductRepo = repository.NewProductRepository(c.DB)
	c.QuoteRepo = repository.NewQuoteRepository(c.DB)
	c.QuoteItemRepo = repository.NewQuoteItemRepository(c.DB)
	c.SettingsRepo = repository.NewSettingsRepository(c.DB)
	log.Println("Repositories initialized successfully")
	return nil
}

func (c *Container) initializeUseCases() error {
	if c.UserRepo == nil {
		return fmt.Errorf("repositories not initialized")
	}

	c.TenantUseCase = tenantUseCase.NewUseCase(c.TenantRepo)
	c.UserUseCase = userUseCase.NewUseCase(c.UserRepo, c.JWTManager, c.PassHasher)
	c.ClientUseCase = clientUseCase.NewUseCase(c.ClientRepo)
	c.ProductUseCase = productUseCase.NewUseCase(c.ProductRepo)
	c.QuoteUseCase = quoteUseCase.NewUseCase(c.QuoteRepo, c.QuoteItemRepo)
	c.SettingsUseCase = settingsUseCase.NewUseCase(c.SettingsRepo)
	log.Printf("Use cases initialized successfully - TenantRepo: %v, UserRepo: %v, ClientRepo: %v, ProductRepo: %v, QuoteRepo: %v, SettingsRepo: %v, JWTManager: %v, PassHasher: %v",
		c.TenantRepo != nil, c.UserRepo != nil, c.ClientRepo != nil, c.ProductRepo != nil, c.QuoteRepo != nil, c.SettingsRepo != nil, c.JWTManager != nil, c.PassHasher != nil)
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
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (c *Container) GetUserRepository() userDomain.Repository {
	return c.UserRepo
}

func (c *Container) GetUserUseCase() userUseCase.UseCaseInterface {
	return c.UserUseCase
}

func (c *Container) GetClientRepository() clientDomain.Repository {
	return c.ClientRepo
}

func (c *Container) GetClientUseCase() clientUseCase.UseCaseInterface {
	return c.ClientUseCase
}

func (c *Container) GetProductRepository() productDomain.Repository {
	return c.ProductRepo
}

func (c *Container) GetProductUseCase() productUseCase.UseCaseInterface {
	return c.ProductUseCase
}

func (c *Container) GetQuoteRepository() quoteDomain.Repository {
	return c.QuoteRepo
}

func (c *Container) GetQuoteItemRepository() quoteDomain.ItemRepository {
	return c.QuoteItemRepo
}

func (c *Container) GetQuoteUseCase() quoteUseCase.UseCaseInterface {
	return c.QuoteUseCase
}

func (c *Container) GetSettingsRepository() settingsDomain.Repository {
	return c.SettingsRepo
}

func (c *Container) GetSettingsUseCase() settingsUseCase.UseCaseInterface {
	return c.SettingsUseCase
}

func (c *Container) GetTenantRepository() tenantDomain.Repository {
	return c.TenantRepo
}

func (c *Container) GetTenantUseCase() tenantUseCase.UseCaseInterface {
	return c.TenantUseCase
}

func (c *Container) GetJWTManager() *auth.JWTManager {
	return c.JWTManager
}

func (c *Container) GetPassHasher() *auth.PasswordHasher {
	return c.PassHasher
} 