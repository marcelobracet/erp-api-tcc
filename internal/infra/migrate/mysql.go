package migrate

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	auditDomain "erp-api/internal/domain/audit"
	clientDomain "erp-api/internal/domain/client"
	productDomain "erp-api/internal/domain/product"
	quoteDomain "erp-api/internal/domain/quote"
	settingsDomain "erp-api/internal/domain/settings"
	tenantDomain "erp-api/internal/domain/tenant"
	userDomain "erp-api/internal/domain/user"

	"gorm.io/gorm"
)

type MySQLMigrator struct{}

func (m *MySQLMigrator) Run(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Running database migrations (mysql)...")

	if err := ensureTenantCompanyNameColumnMySQL(db); err != nil {
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

	createForeignKeysMySQL(db)
	createGeneratedColumnsAndIndexesMySQL(db)

	log.Println("Database migrations completed successfully (mysql)")
	return nil
}

func ensureTenantCompanyNameColumnMySQL(db *gorm.DB) error {
	// only run if tenants table exists (first deploy won't have anything yet)
	tableExists, err := mysqlTableExists(db, "tenants")
	if err != nil {
		return err
	}
	if !tableExists {
		return nil
	}

	companyNameExists, err := mysqlColumnExists(db, "tenants", "company_name")
	if err != nil {
		return err
	}
	if !companyNameExists {
		if err := db.Exec("ALTER TABLE tenants ADD COLUMN company_name TEXT NULL").Error; err != nil {
			return err
		}
	}

	hasTradeName, err := mysqlColumnExists(db, "tenants", "trade_name")
	if err != nil {
		return err
	}
	hasEmail, err := mysqlColumnExists(db, "tenants", "email")
	if err != nil {
		return err
	}

	// backfill null/empty values
	var update string
	switch {
	case hasTradeName && hasEmail:
		update = "UPDATE tenants SET company_name = COALESCE(NULLIF(trade_name,''), NULLIF(email,''), 'Legacy Tenant') WHERE company_name IS NULL OR company_name = ''"
	case hasTradeName:
		update = "UPDATE tenants SET company_name = COALESCE(NULLIF(trade_name,''), 'Legacy Tenant') WHERE company_name IS NULL OR company_name = ''"
	case hasEmail:
		update = "UPDATE tenants SET company_name = COALESCE(NULLIF(email,''), 'Legacy Tenant') WHERE company_name IS NULL OR company_name = ''"
	default:
		update = "UPDATE tenants SET company_name = 'Legacy Tenant' WHERE company_name IS NULL OR company_name = ''"
	}
	if err := db.Exec(update).Error; err != nil {
		return err
	}

	// enforce NOT NULL
	if err := db.Exec("ALTER TABLE tenants MODIFY company_name TEXT NOT NULL").Error; err != nil {
		return err
	}

	return nil
}

func createForeignKeysMySQL(db *gorm.DB) {
	// Best-effort: if it fails due to types/indexes, we log and continue
	addFKIfMissing(db, "users", "fk_users_tenant", "ALTER TABLE users ADD CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE")
	addFKIfMissing(db, "audits", "fk_audits_tenant", "ALTER TABLE audits ADD CONSTRAINT fk_audits_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE")
	addFKIfMissing(db, "audits", "fk_audits_user", "ALTER TABLE audits ADD CONSTRAINT fk_audits_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL")
	addFKIfMissing(db, "clients", "fk_clients_tenant", "ALTER TABLE clients ADD CONSTRAINT fk_clients_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE")
	addFKIfMissing(db, "products", "fk_products_tenant", "ALTER TABLE products ADD CONSTRAINT fk_products_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE")

	addFKIfMissing(db, "quotes", "fk_quotes_tenant", "ALTER TABLE quotes ADD CONSTRAINT fk_quotes_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE")
	addFKIfMissing(db, "quotes", "fk_quotes_client", "ALTER TABLE quotes ADD CONSTRAINT fk_quotes_client FOREIGN KEY (client_id) REFERENCES clients(id)")
	addFKIfMissing(db, "quotes", "fk_quotes_user", "ALTER TABLE quotes ADD CONSTRAINT fk_quotes_user FOREIGN KEY (user_id) REFERENCES users(id)")

	addFKIfMissing(db, "quote_items", "fk_quote_items_tenant", "ALTER TABLE quote_items ADD CONSTRAINT fk_quote_items_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE")
	addFKIfMissing(db, "quote_items", "fk_quote_items_quote", "ALTER TABLE quote_items ADD CONSTRAINT fk_quote_items_quote FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE CASCADE")
	addFKIfMissing(db, "quote_items", "fk_quote_items_product", "ALTER TABLE quote_items ADD CONSTRAINT fk_quote_items_product FOREIGN KEY (product_id) REFERENCES products(id)")

	addFKIfMissing(db, "settings", "fk_settings_tenant", "ALTER TABLE settings ADD CONSTRAINT fk_settings_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE")
}

func createGeneratedColumnsAndIndexesMySQL(db *gorm.DB) {
	// generated column: quote_items.total
	totalExists, err := mysqlColumnExists(db, "quote_items", "total")
	if err != nil {
		log.Printf("Warning: could not check quote_items.total: %v", err)
	} else if !totalExists {
		if err := db.Exec("ALTER TABLE quote_items ADD COLUMN total DECIMAL(12,2) GENERATED ALWAYS AS (quantity * price) STORED").Error; err != nil {
			log.Printf("Warning: could not create generated column quote_items.total: %v", err)
		}
	}

	// unique indexes
	createIndexIfMissing(db, "clients", "idx_clients_document_tenant", "CREATE UNIQUE INDEX idx_clients_document_tenant ON clients(document, tenant_id)")
	createIndexIfMissing(db, "settings", "idx_settings_key_tenant", "CREATE UNIQUE INDEX idx_settings_key_tenant ON settings(`key`, tenant_id)")

	// CHECK constraint support is version-dependent in MySQL. Best-effort only.
	if supportsMySQLCheckConstraints(db) {
		addCheckIfPossible(db, "clients", "clients_document_type_check", "ALTER TABLE clients ADD CONSTRAINT clients_document_type_check CHECK (document_type IN ('CPF','CNPJ'))")
	} else {
		log.Printf("Info: skipping CHECK constraint clients_document_type_check (MySQL version without enforced CHECK)")
	}
}

func mysqlTableExists(db *gorm.DB, table string) (bool, error) {
	var count int
	err := db.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`, table).Scan(&count).Error
	return count > 0, err
}

func mysqlColumnExists(db *gorm.DB, table, column string) (bool, error) {
	var count int
	err := db.Raw(`SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?`, table, column).Scan(&count).Error
	return count > 0, err
}

func mysqlConstraintExists(db *gorm.DB, table, constraint string) (bool, error) {
	var count int
	err := db.Raw(`SELECT COUNT(*) FROM information_schema.table_constraints WHERE constraint_schema = DATABASE() AND table_name = ? AND constraint_name = ?`, table, constraint).Scan(&count).Error
	return count > 0, err
}

func mysqlIndexExists(db *gorm.DB, table, index string) (bool, error) {
	var count int
	err := db.Raw(`SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?`, table, index).Scan(&count).Error
	return count > 0, err
}

func addFKIfMissing(db *gorm.DB, table, constraint, sqlStmt string) {
	exists, err := mysqlConstraintExists(db, table, constraint)
	if err != nil {
		log.Printf("Warning: could not check FK %s on %s: %v", constraint, table, err)
		return
	}
	if exists {
		return
	}
	if err := db.Exec(sqlStmt).Error; err != nil {
		log.Printf("Warning: could not create FK %s on %s: %v", constraint, table, err)
	}
}

func createIndexIfMissing(db *gorm.DB, table, index, sqlStmt string) {
	exists, err := mysqlIndexExists(db, table, index)
	if err != nil {
		log.Printf("Warning: could not check index %s on %s: %v", index, table, err)
		return
	}
	if exists {
		return
	}
	if err := db.Exec(sqlStmt).Error; err != nil {
		log.Printf("Warning: could not create index %s on %s: %v", index, table, err)
	}
}

func addCheckIfPossible(db *gorm.DB, table, constraint, sqlStmt string) {
	exists, err := mysqlConstraintExists(db, table, constraint)
	if err != nil {
		log.Printf("Warning: could not check constraint %s on %s: %v", constraint, table, err)
		return
	}
	if exists {
		return
	}
	if err := db.Exec(sqlStmt).Error; err != nil {
		log.Printf("Warning: could not create CHECK %s on %s: %v", constraint, table, err)
	}
}

func supportsMySQLCheckConstraints(db *gorm.DB) bool {
	// MySQL enforces CHECK constraints since 8.0.16.
	var versionString string
	if err := db.Raw("SELECT VERSION()").Scan(&versionString).Error; err != nil {
		log.Printf("Warning: could not detect MySQL version: %v", err)
		return false
	}
	major, minor, patch := parseMySQLVersion(versionString)
	if major == 0 {
		return false
	}
	if major > 8 {
		return true
	}
	if major < 8 {
		return false
	}
	if minor > 0 {
		return true
	}
	return patch >= 16
}

func parseMySQLVersion(v string) (major, minor, patch int) {
	// e.g. "8.0.35" or "8.0.35-log"
	clean := v
	for i, r := range v {
		if !(r >= '0' && r <= '9') && r != '.' {
			clean = v[:i]
			break
		}
	}
	parts := strings.Split(clean, ".")
	if len(parts) < 2 {
		return 0, 0, 0
	}
	major, _ = strconv.Atoi(parts[0])
	minor, _ = strconv.Atoi(parts[1])
	if len(parts) >= 3 {
		patch, _ = strconv.Atoi(parts[2])
	}
	return major, minor, patch
}
