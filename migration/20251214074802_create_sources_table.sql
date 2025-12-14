-- +goose Up
-- +goose StatementBegin
CREATE TYPE source_type AS ENUM (
    'link',
    'pdf',
    'ppt',
    'doc',
    'note'
);

CREATE TYPE source_status AS ENUM (
    'pending',
    'processing',
    'indexed',
    'failed'
);

CREATE TABLE IF NOT EXISTS sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    collection_id UUID REFERENCES collections(id) ON DELETE SET NULL,

    type source_type NOT NULL,
    status source_status NOT NULL DEFAULT 'pending',

    title TEXT NOT NULL,

    -- For links
    original_url TEXT,

    -- For uploaded files
    s3_bucket TEXT,
    s3_key TEXT,

    -- Dedup / reprocessing safety
    content_hash TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sources_user_id ON sources(user_id);
CREATE INDEX IF NOT EXISTS idx_sources_collection_id ON sources(collection_id);
CREATE INDEX IF NOT EXISTS idx_sources_status ON sources(status);
CREATE INDEX IF NOT EXISTS idx_sources_type ON sources(type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sources_user_id;
DROP INDEX IF EXISTS idx_sources_collection_id;
DROP INDEX IF EXISTS idx_sources_status;
DROP INDEX IF EXISTS idx_sources_type;
DROP TABLE IF EXISTS sources;
-- +goose StatementEnd

