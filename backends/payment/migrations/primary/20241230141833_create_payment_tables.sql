-- migrate:up
CREATE TYPE payment_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'canceled',
    'refunded'
);

CREATE TYPE payment_provider AS ENUM (
    'stripe',
    'adyen'
);

CREATE TYPE payment_method_type AS ENUM (
    'card',
    'bank_transfer',
    'crypto'
);

CREATE TABLE payments (
    id TEXT PRIMARY KEY DEFAULT uuid7(),
    merchant_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status payment_status NOT NULL,
    provider payment_provider NOT NULL,
    method payment_method_type NOT NULL,
    description TEXT,
    provider_id VARCHAR(255),
    error_message TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_payments_merchant_id ON payments(merchant_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);

CREATE TABLE payment_intents (
    id TEXT PRIMARY KEY DEFAULT uuid7(),
    merchant_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status payment_status NOT NULL,
    provider payment_provider NOT NULL,
    method payment_method_type NOT NULL,
    description TEXT,
    client_secret VARCHAR(255) NOT NULL,
    return_url TEXT NOT NULL,
    webhook_url TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_payment_intents_merchant_id ON payment_intents(merchant_id);
CREATE INDEX idx_payment_intents_status ON payment_intents(status);
CREATE INDEX idx_payment_intents_expires_at ON payment_intents(expires_at);

CREATE TABLE payment_methods (
    id TEXT PRIMARY KEY DEFAULT uuid7(),
    merchant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    type payment_method_type NOT NULL,
    provider payment_provider NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    last4 VARCHAR(4),
    expiry_month INTEGER,
    expiry_year INTEGER,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_payment_methods_customer_id ON payment_methods(customer_id);
CREATE INDEX idx_payment_methods_merchant_id ON payment_methods(merchant_id);

-- migrate:down
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_intents;
DROP TABLE IF EXISTS payment_methods;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_provider;
DROP TYPE IF EXISTS payment_method_type;
