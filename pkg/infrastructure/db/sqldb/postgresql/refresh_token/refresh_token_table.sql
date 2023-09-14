DROP TABLE IF EXISTS refresh_token;
CREATE TABLE IF NOT EXISTS refresh_token
(
    id          SERIAL PRIMARY KEY,
    app_user_id BIGINT      NOT NULL references app_user (id) ON DELETE CASCADE,
    token_uuid  UUID        NOT NULL,
    revoked     BOOLEAN     NOT NULL DEFAULT FALSE,
    expires_at  TIMESTAMPTZ NOT NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_token_uuid ON refresh_token (token_uuid);
