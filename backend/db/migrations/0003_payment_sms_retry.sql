BEGIN;

CREATE TABLE IF NOT EXISTS payment_attempts (
    id BIGSERIAL PRIMARY KEY,
    service_request_id BIGINT NOT NULL REFERENCES service_requests(id),
    gateway_order_id VARCHAR(120) NOT NULL,
    gateway_payment_id VARCHAR(120),
    amount NUMERIC(10,2) NOT NULL CHECK (amount >= 0.01),
    status VARCHAR(30) NOT NULL CHECK (status IN ('created', 'authorized', 'captured', 'failed')),
    idempotency_key VARCHAR(128) NOT NULL,
    callback_payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT payment_attempts_gateway_order_id_unique UNIQUE (gateway_order_id),
    CONSTRAINT payment_attempts_gateway_payment_id_unique UNIQUE (gateway_payment_id)
);

CREATE TABLE IF NOT EXISTS service_request_retry_audit (
    id BIGSERIAL PRIMARY KEY,
    service_request_id BIGINT NOT NULL REFERENCES service_requests(id),
    retry_no INTEGER NOT NULL CHECK (retry_no >= 1),
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sms_notifications (
    id BIGSERIAL PRIMARY KEY,
    service_request_id BIGINT NOT NULL REFERENCES service_requests(id),
    phone_number VARCHAR(10) NOT NULL,
    template TEXT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('queued', 'sent', 'failed')),
    attempt_count INTEGER NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
    last_error TEXT,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_payment_attempts_service_request_idempotency
    ON payment_attempts (service_request_id, idempotency_key);

CREATE INDEX IF NOT EXISTS idx_payment_attempts_service_request
    ON payment_attempts (service_request_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_retry_audit_service_request
    ON service_request_retry_audit (service_request_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_sms_notifications_service_request
    ON sms_notifications (service_request_id, created_at DESC);

COMMIT;
