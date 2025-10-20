CREATE TABLE IF NOT EXISTS quotes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    client_id UUID NOT NULL,
    user_id UUID NOT NULL,
    total_value NUMERIC(12,2) DEFAULT 0,
    discount NUMERIC(10,2) DEFAULT 0,
    status VARCHAR(50) DEFAULT 'pending',
    conversion_rate NUMERIC(5,2),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_quotes_tenant_id ON quotes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_quotes_status ON quotes(status);
