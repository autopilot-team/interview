-- migrate:up

-- Create enum for payment status
CREATE TYPE payment_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'cancelled',
    'partially_refunded',
    'refunded'
);

-- Create enum for payment method types
CREATE TYPE payment_method_type AS ENUM (
    'card',
    'bank_transfer',
    'wallet',
    'other'
);

-- Create enum for refund status
CREATE TYPE refund_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'cancelled'
);

-- Create enum for refund reason
CREATE TYPE refund_reason AS ENUM (
    'duplicate',
    'fraudulent',
    'customer_request',
    'product_issue',
    'other'
);

-- Create merchants table
CREATE TABLE IF NOT EXISTS merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    payment_provider VARCHAR(50) NOT NULL, -- 'stripe', 'adyen', etc.
    provider_merchant_id VARCHAR(255) NOT NULL, -- External merchant ID from payment provider
    settings JSONB NOT NULL DEFAULT '{}', -- Provider-specific settings
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT merchants_entity_id_fkey FOREIGN KEY (entity_id) REFERENCES identity.entities(id) ON DELETE CASCADE,
    UNIQUE (payment_provider, provider_merchant_id)
);

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL, -- Prevents duplicate payments
    external_payment_id VARCHAR(255), -- Payment ID from payment provider
    amount BIGINT NOT NULL, -- Amount in smallest currency unit (e.g., cents)
    currency VARCHAR(3) NOT NULL, -- ISO 4217 currency code
    status payment_status NOT NULL DEFAULT 'pending',
    payment_method payment_method_type NOT NULL,
    payment_method_details JSONB NOT NULL DEFAULT '{}', -- Card last 4, bank name, etc.
    customer_id UUID, -- Optional reference to customer in identity system
    customer_email VARCHAR(255),
    customer_name VARCHAR(255),
    description TEXT,
    metadata JSONB NOT NULL DEFAULT '{}', -- Custom key-value pairs
    
    -- Payment provider response data
    provider_response JSONB NOT NULL DEFAULT '{}',
    provider_error_code VARCHAR(100),
    provider_error_message TEXT,
    
    -- Refund tracking
    refunded_amount BIGINT NOT NULL DEFAULT 0,
    refundable_amount BIGINT NOT NULL DEFAULT 0,
    refund_count INTEGER NOT NULL DEFAULT 0,
    
    -- Important timestamps
    initiated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT payments_merchant_id_fkey FOREIGN KEY (merchant_id) REFERENCES merchants(id) ON DELETE RESTRICT,
    CONSTRAINT payments_amount_positive CHECK (amount > 0),
    CONSTRAINT payments_refunded_amount_check CHECK (refunded_amount >= 0 AND refunded_amount <= amount),
    UNIQUE (merchant_id, idempotency_key)
);

-- Create refunds table
CREATE TABLE IF NOT EXISTS refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL, -- Prevents duplicate refunds
    external_refund_id VARCHAR(255), -- Refund ID from payment provider
    amount BIGINT NOT NULL, -- Amount to refund in smallest currency unit
    currency VARCHAR(3) NOT NULL, -- Must match payment currency
    status refund_status NOT NULL DEFAULT 'pending',
    reason refund_reason NOT NULL,
    reason_description TEXT,
    metadata JSONB NOT NULL DEFAULT '{}',
    
    -- Payment provider response data
    provider_response JSONB NOT NULL DEFAULT '{}',
    provider_error_code VARCHAR(100),
    provider_error_message TEXT,
    
    -- User who initiated the refund
    initiated_by UUID, -- User ID from identity system
    initiated_by_email VARCHAR(255),
    
    -- Important timestamps
    initiated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT refunds_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE RESTRICT,
    CONSTRAINT refunds_amount_positive CHECK (amount > 0),
    UNIQUE (payment_id, idempotency_key)
);

-- Create payment_events table for audit trail
CREATE TABLE IF NOT EXISTS payment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- 'created', 'processing', 'succeeded', 'failed', 'refund_initiated', etc.
    event_data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT payment_events_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_merchants_entity_id ON merchants(entity_id);
CREATE INDEX idx_merchants_is_active ON merchants(is_active);

CREATE INDEX idx_payments_merchant_id ON payments(merchant_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_customer_id ON payments(customer_id);
CREATE INDEX idx_payments_customer_email ON payments(customer_email);
CREATE INDEX idx_payments_external_payment_id ON payments(external_payment_id);
CREATE INDEX idx_payments_created_at ON payments(created_at DESC);

CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_status ON refunds(status);
CREATE INDEX idx_refunds_initiated_by ON refunds(initiated_by);
CREATE INDEX idx_refunds_created_at ON refunds(created_at DESC);

CREATE INDEX idx_payment_events_payment_id ON payment_events(payment_id);
CREATE INDEX idx_payment_events_created_at ON payment_events(created_at DESC);

-- Create triggers to update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_merchants_updated_at BEFORE UPDATE ON merchants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_refunds_updated_at BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- migrate:down

DROP TRIGGER IF EXISTS update_refunds_updated_at ON refunds;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_merchants_updated_at ON merchants;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS payment_events;
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS merchants;

DROP TYPE IF EXISTS refund_reason;
DROP TYPE IF EXISTS refund_status;
DROP TYPE IF EXISTS payment_method_type;
DROP TYPE IF EXISTS payment_status;