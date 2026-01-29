-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS source_contents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    content_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_source_contents_source_id ON source_contents(source_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_source_contents_source_id;
DROP TABLE IF EXISTS source_contents;
-- +goose StatementEnd
