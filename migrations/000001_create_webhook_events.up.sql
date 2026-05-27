CREATE TABLE webhook_events (
    id TEXT PRIMARY KEY,
    target_url TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_webhook_events_created_at
    ON webhook_events (created_at);
