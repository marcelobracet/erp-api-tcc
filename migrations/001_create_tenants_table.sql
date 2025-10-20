CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    plan VARCHAR(50) DEFAULT 'free',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tenants_name ON tenants(name);
CREATE INDEX IF NOT EXISTS idx_tenants_plan ON tenants(plan);
CREATE INDEX IF NOT EXISTS idx_tenants_is_active ON tenants(is_active);
CREATE INDEX IF NOT EXISTS idx_tenants_created_at ON tenants(created_at);
CREATE INDEX IF NOT EXISTS idx_tenants_updated_at ON tenants(updated_at);