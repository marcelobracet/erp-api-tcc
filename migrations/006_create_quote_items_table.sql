CREATE TABLE IF NOT EXISTS quote_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    quote_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    price NUMERIC(10,2) NOT NULL,
    total NUMERIC(12,2) GENERATED ALWAYS AS (quantity * price) STORED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX IF NOT EXISTS idx_quote_items_tenant_id ON quote_items(tenant_id);
