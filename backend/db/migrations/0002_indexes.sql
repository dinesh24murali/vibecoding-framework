BEGIN;

CREATE UNIQUE INDEX IF NOT EXISTS idx_providers_name_active_unique
    ON providers (LOWER(name))
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_plans_provider_name_active_unique
    ON plans (provider_id, LOWER(name))
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_service_requests_status_created_at
    ON service_requests (status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_service_requests_user_created_at
    ON service_requests (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_plans_provider_id
    ON plans (provider_id);

CREATE INDEX IF NOT EXISTS idx_users_role_status
    ON users (role, status);

COMMIT;
