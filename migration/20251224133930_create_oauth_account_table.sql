-- +goose Up
-- +goose StatementBegin
------------------------------------------------
-- OAUTH ACCOUNTS
------------------------------------------------
CREATE TYPE oauth_provider AS ENUM ('google', 'github');

CREATE TABLE IF NOT EXISTS oauth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider oauth_provider NOT NULL,
    provider_user_id TEXT NOT NULL,
    email TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(provider, provider_user_id)
);

CREATE INDEX IF NOT EXISTS idx_oauth_accounts_user_id ON oauth_accounts(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_oauth_accounts_user_id;
DROP TABLE IF EXISTS oauth_accounts;
DROP TYPE IF EXISTS oauth_provider;
-- +goose StatementEnd
