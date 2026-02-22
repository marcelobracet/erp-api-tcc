package migrate

import (
	"fmt"
	"log"

	auditDomain "erp-api/internal/domain/audit"
	clientDomain "erp-api/internal/domain/client"
	productDomain "erp-api/internal/domain/product"
	quoteDomain "erp-api/internal/domain/quote"
	settingsDomain "erp-api/internal/domain/settings"
	tenantDomain "erp-api/internal/domain/tenant"
	userDomain "erp-api/internal/domain/user"

	"gorm.io/gorm"
)

type PostgresMigrator struct{}

func (m *PostgresMigrator) Run(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Running database migrations (postgres)...")

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: Could not create uuid-ossp extension: %v", err)
		if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
			log.Printf("Warning: Could not create pgcrypto extension: %v", err)
		}
	}

	if err := ensureTenantCompanyNameColumnPostgres(db); err != nil {
		return fmt.Errorf("failed to ensure tenants.company_name: %w", err)
	}

	if err := db.AutoMigrate(
		&auditDomain.Audit{},
		&tenantDomain.Tenant{},
		&userDomain.User{},
		&clientDomain.Client{},
		&productDomain.Product{},
		&quoteDomain.Quote{},
		&quoteDomain.QuoteItem{},
		&settingsDomain.Settings{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	createForeignKeysPostgres(db)
	createGeneratedColumnsAndIndexesPostgres(db)

	log.Println("Database migrations completed successfully (postgres)")
	return nil
}

func ensureTenantCompanyNameColumnPostgres(db *gorm.DB) error {
	res := db.Exec(`
		DO $$
		DECLARE
			has_trade_name boolean;
			has_email boolean;
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = 'tenants'
			) THEN
				IF NOT EXISTS (
					SELECT 1
					FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'tenants' AND column_name = 'company_name'
				) THEN
					ALTER TABLE tenants ADD COLUMN company_name text;
				END IF;

				SELECT EXISTS (
					SELECT 1
					FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'tenants' AND column_name = 'trade_name'
				) INTO has_trade_name;

				SELECT EXISTS (
					SELECT 1
					FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'tenants' AND column_name = 'email'
				) INTO has_email;

				IF has_trade_name AND has_email THEN
					EXECUTE 'UPDATE tenants SET company_name = COALESCE(NULLIF(trade_name, ''''), NULLIF(email, ''''), ''Legacy Tenant'') WHERE company_name IS NULL OR company_name = ''''';
				ELSIF has_trade_name THEN
					EXECUTE 'UPDATE tenants SET company_name = COALESCE(NULLIF(trade_name, ''''), ''Legacy Tenant'') WHERE company_name IS NULL OR company_name = ''''';
				ELSIF has_email THEN
					EXECUTE 'UPDATE tenants SET company_name = COALESCE(NULLIF(email, ''''), ''Legacy Tenant'') WHERE company_name IS NULL OR company_name = ''''';
				ELSE
					EXECUTE 'UPDATE tenants SET company_name = ''Legacy Tenant'' WHERE company_name IS NULL OR company_name = ''''';
				END IF;

				ALTER TABLE tenants ALTER COLUMN company_name SET NOT NULL;
			END IF;
		END $$;
	`)
	return res.Error
}

func createForeignKeysPostgres(db *gorm.DB) {
	db.Exec(`
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

	db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_audits_tenant'
			) THEN
				ALTER TABLE audits ADD CONSTRAINT fk_audits_tenant 
				FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`)

	db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_audits_user'
			) THEN
				ALTER TABLE audits ADD CONSTRAINT fk_audits_user 
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
			END IF;
		END $$;
	`)

	db.Exec(`
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

	db.Exec(`
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

	db.Exec(`
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

	db.Exec(`
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

	db.Exec(`
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

func createGeneratedColumnsAndIndexesPostgres(db *gorm.DB) {
	db.Exec(`
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

	db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_indexes WHERE indexname = 'idx_clients_document_tenant'
			) THEN
				CREATE UNIQUE INDEX idx_clients_document_tenant ON clients(document, tenant_id);
			END IF;
		END $$;
	`)

	db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_indexes WHERE indexname = 'idx_settings_key_tenant'
			) THEN
				CREATE UNIQUE INDEX idx_settings_key_tenant ON settings(key, tenant_id);
			END IF;
		END $$;
	`)

	db.Exec(`
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
