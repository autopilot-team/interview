-- migrate:up
CREATE TABLE "payments" (
    "id" UUID NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "merchant_id" UUID NOT NULL, -- references entities(id)? or create a new merchants table?
    "amount" BIGINT NOT NULL,
    "currency" TEXT NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "provider" TEXT NOT NULL,
    "method" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "error_message" TEXT,
    "metadata" JSONB DEFAULT '{}',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "completed_at" TIMESTAMPTZ,
    CONSTRAINT "valid_payment_status" CHECK (status IN ('pending', 'processing', 'succeeded', 'failed', 'canceled', 'refunded')),
    CONSTRAINT "valid_payment_method" CHECK (method IN ('card', 'bank_transfer'))
);

CREATE INDEX idx_payments_merchant_id ON payments(merchant_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);

COMMENT ON TABLE "payments" IS 'Manage payment transactions.';

CREATE TABLE "payment_intents" (
    "id" UUID NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "merchant_id" UUID NOT NULL, -- references entities(id)? or create a new merchants table?
    "amount" BIGINT NOT NULL,
    "currency" TEXT NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "method" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "client_secret" TEXT NOT NULL,
    "return_url" TEXT NOT NULL,
    "webhook_url" TEXT NOT NULL,
    "metadata" JSONB DEFAULT '{}',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "expires_at" TIMESTAMPTZ NOT NULL,
    CONSTRAINT "valid_payment_intent_status" CHECK (status IN ('pending', 'processing', 'succeeded', 'failed', 'canceled', 'refunded')),
    CONSTRAINT "valid_payment_intent_method" CHECK (method IN ('card', 'bank_transfer'))
);

CREATE INDEX idx_payment_intents_merchant_id ON payment_intents(merchant_id);
CREATE INDEX idx_payment_intents_status ON payment_intents(status);
CREATE INDEX idx_payment_intents_expires_at ON payment_intents(expires_at);

COMMENT ON TABLE "payment_intents" IS 'Manage payment intents before processing.';

CREATE TABLE "stored_payment_methods" (
    "id" UUID NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "merchant_id" UUID NOT NULL,
    "customer_id" UUID NOT NULL,
    "type" TEXT NOT NULL,
    "provider_id" TEXT NOT NULL,
    "last4" VARCHAR(4) NOT NULL,
    "expiry_month" INTEGER NOT NULL,
    "expiry_year" INTEGER NOT NULL,
    "metadata" JSONB DEFAULT '{}',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "valid_stored_payment_method_type" CHECK (type IN ('card', 'bank_transfer'))
);

CREATE INDEX idx_stored_payment_methods_merchant_id ON stored_payment_methods(merchant_id);
CREATE INDEX idx_stored_payment_methods_customer_id ON stored_payment_methods(customer_id);
CREATE INDEX idx_stored_payment_methods_deleted_at ON stored_payment_methods(deleted_at);

COMMENT ON TABLE "stored_payment_methods" IS 'Manage stored customer payment methods.';

-- migrate:down
DROP TABLE "stored_payment_methods";
DROP TABLE "payment_intents";
DROP TABLE "payments";
