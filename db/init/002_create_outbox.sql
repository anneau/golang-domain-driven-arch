CREATE TABLE IF NOT EXISTS outbox (
    id           TEXT PRIMARY KEY,
    event_type   TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    payload      JSONB NOT NULL,
    occurred_at  TIMESTAMPTZ NOT NULL,
    published_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS outbox_unpublished_idx ON outbox (occurred_at) WHERE published_at IS NULL;
