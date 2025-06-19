-- Script to seed test merchant data

-- First, ensure we have a test entity in identity database
\c identity

INSERT INTO entities (id, name, slug, type, status) 
VALUES (
    '550e8400-e29b-41d4-a716-446655440000'::uuid,
    'Test Organization',
    'test-org',
    'organization',
    'active'
) ON CONFLICT (id) DO NOTHING;

-- Add merchant to live database
\c payment_live

INSERT INTO merchants (
    id, 
    entity_id, 
    name, 
    payment_provider, 
    provider_merchant_id, 
    is_active,
    settings
) VALUES (
    '550e8400-e29b-41d4-a716-446655440001'::uuid,
    '550e8400-e29b-41d4-a716-446655440000'::uuid,
    'Test Merchant - Live',
    'stripe',
    'stripe_live_merchant_123',
    true,
    '{"api_key": "live_test_key", "webhook_secret": "live_test_secret"}'::jsonb
) ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    is_active = EXCLUDED.is_active;

-- Add merchant to test database
\c payment_test

INSERT INTO merchants (
    id, 
    entity_id, 
    name, 
    payment_provider, 
    provider_merchant_id, 
    is_active,
    settings
) VALUES (
    '550e8400-e29b-41d4-a716-446655440001'::uuid,
    '550e8400-e29b-41d4-a716-446655440000'::uuid,
    'Test Merchant - Test',
    'stripe',
    'stripe_test_merchant_123',
    true,
    '{"api_key": "test_test_key", "webhook_secret": "test_test_secret"}'::jsonb
) ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    is_active = EXCLUDED.is_active;

-- Verify the data
SELECT * FROM merchants;