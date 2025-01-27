-- migrate:up
CREATE TABLE "users" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL UNIQUE,
    "email_verified_at" TIMESTAMPTZ NULL,
    "failed_login_attempts" INTEGER NOT NULL DEFAULT 0,
    "image" TEXT,
    "last_active_at" TIMESTAMPTZ,
    "last_logged_in_at" TIMESTAMPTZ,
    "locked_at" TIMESTAMPTZ,
    "password_changed_at" TIMESTAMPTZ,
    "password_hash" TEXT,
    "two_factor_enabled" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE "users" IS 'Manage user information.';

CREATE TABLE "platforms" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "domain" TEXT,
    "logo" TEXT,
    "name" TEXT NOT NULL,
    "slug" TEXT UNIQUE,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT "valid_platform_status" CHECK (status IN ('pending', 'active', 'inactive', 'suspended'))
);

COMMENT ON TABLE "platforms" IS 'Manage platforms.';

CREATE TABLE "organizations" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "domain" TEXT,
    "logo" TEXT,
    "name" TEXT NOT NULL,
    "platform_id" TEXT REFERENCES "platforms" ("id") ON DELETE CASCADE,
    "slug" TEXT UNIQUE,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT "valid_organization_status" CHECK (status IN ('pending', 'active', 'inactive', 'suspended'))
);
CREATE INDEX idx_organizations_platform_id ON organizations(platform_id);

COMMENT ON TABLE "organizations" IS 'Manage organizations.';

CREATE TABLE "accounts" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "domain" TEXT,
    "logo" TEXT,
    "name" TEXT NOT NULL,
    "organization_id" TEXT REFERENCES "organizations" ("id") ON DELETE CASCADE,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT "valid_account_status" CHECK (status IN ('pending', 'active', 'inactive', 'suspended'))
);
CREATE INDEX idx_accounts_organization_id ON accounts(organization_id);

COMMENT ON TABLE "accounts" IS 'Manage business accounts.';

CREATE TABLE "members" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "account_id" TEXT REFERENCES "accounts" ("id") ON DELETE CASCADE,
    "organization_id" TEXT REFERENCES "organizations" ("id") ON DELETE CASCADE,
    "platform_id" TEXT REFERENCES "platforms" ("id") ON DELETE CASCADE,
    "role" TEXT NOT NULL,
    "user_id" TEXT NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT member_entity_check CHECK (
        (CASE WHEN account_id IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN organization_id IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN platform_id IS NOT NULL THEN 1 ELSE 0 END) = 1
    )
);
CREATE INDEX idx_members_user_id ON members(user_id);
CREATE INDEX idx_members_platform_id ON members(platform_id);
CREATE INDEX idx_members_organization_id ON members(organization_id);
CREATE INDEX idx_members_account_id ON members(account_id);

COMMENT ON TABLE "members" IS 'Manage account/organization/platform members.';

CREATE TABLE "invitations" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "email" TEXT NOT NULL,
    "expires_at" TIMESTAMPTZ NOT NULL,
    "account_id" TEXT REFERENCES "accounts" ("id") ON DELETE CASCADE,
    "organization_id" TEXT REFERENCES "organizations" ("id") ON DELETE CASCADE,
    "platform_id" TEXT REFERENCES "platforms" ("id") ON DELETE CASCADE,
    "inviter_id" TEXT NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "role" TEXT,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT "valid_invitation_status" CHECK (status IN ('pending', 'accepted', 'rejected', 'expired')),
    CONSTRAINT "invitation_entity_check" CHECK (
        (CASE WHEN account_id IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN organization_id IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN platform_id IS NOT NULL THEN 1 ELSE 0 END) = 1
    )
);
CREATE INDEX idx_invitations_email ON invitations(email);
CREATE INDEX idx_invitations_status ON invitations(status);

COMMENT ON TABLE "invitations" IS 'Manage account/organization/platform invitations.';

CREATE TABLE "sessions" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "active_account_id" TEXT REFERENCES accounts(id),
    "active_organization_id" TEXT REFERENCES organizations(id),
    "active_platform_id" TEXT REFERENCES platforms(id),
    "expires_at" TIMESTAMPTZ NOT NULL,
    "ip_address" TEXT,
    "token" TEXT NOT NULL UNIQUE,
    "user_agent" TEXT,
    "user_id" TEXT NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT active_entity_check CHECK (
        (CASE WHEN active_account_id IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN active_organization_id IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN active_platform_id IS NOT NULL THEN 1 ELSE 0 END) <= 1
    )
);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

COMMENT ON TABLE "sessions" IS 'Manage user authentication sessions.';

CREATE TABLE "verifications" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "value" TEXT NOT NULL,
    "expires_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE "verifications" IS 'Manage email and other verification processes.';

CREATE TABLE "passkeys" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "backed_up_at" TIMESTAMPTZ NOT NULL,
    "counter" INTEGER NOT NULL,
    "credential_id" TEXT NOT NULL,
    "device_type" TEXT NOT NULL,
    "name" TEXT,
    "public_key" TEXT NOT NULL,
    "transports" TEXT,
    "user_id" TEXT NOT NULL REFERENCES "users" ("id"),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE "passkeys" IS 'Manage user passkeys.';

CREATE TABLE "two_factors" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "backup_codes" TEXT NOT NULL,
    "secret" TEXT NOT NULL,
    "user_id" TEXT NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE "two_factors" IS 'Manage user two-factor authentication.';

CREATE TABLE "compliance_records" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "entity_type" TEXT NOT NULL,
    "entity_id" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "metadata" JSONB DEFAULT '{}',
    "verified_at" TIMESTAMPTZ,
    "verified_by" TEXT REFERENCES users(id),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_type CHECK (type IN ('kyb', 'kyc', 'aml')),
    CONSTRAINT valid_status CHECK (status IN ('pending', 'in_progress', 'approved', 'rejected'))
);

COMMENT ON TABLE "compliance_records" IS 'Manage compliance records.';

CREATE TABLE "audit_logs" (
    "id" TEXT NOT NULL PRIMARY KEY DEFAULT uuid7(),
    "action" TEXT NOT NULL,
    "entity_type" TEXT NOT NULL,
    "entity_id" TEXT NOT NULL,
    "ip_address" TEXT,
    "metadata" JSONB DEFAULT '{}',
    "user_agent" TEXT,
    "user_id" TEXT NOT NULL REFERENCES "users" ("id"),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX "idx_audit_logs_entity" ON "audit_logs"("entity_type", "entity_id");
CREATE INDEX "idx_audit_logs_user_id" ON "audit_logs"("user_id");
CREATE INDEX "idx_audit_logs_created_at" ON "audit_logs"("created_at");

COMMENT ON TABLE "audit_logs" IS 'Manage audit logs.';

-- migrate:down
DROP TABLE "audit_logs";
DROP TABLE "two_factors";
DROP TABLE "passkeys";
DROP TABLE "sessions";
DROP TABLE "verifications";
DROP TABLE "members";
DROP TABLE "invitations";
DROP TABLE "compliance_records";
DROP TABLE "users";
DROP TABLE "organizations";
DROP TABLE "platforms";
DROP TABLE "accounts";
