CREATE TABLE IF NOT EXISTS seed_requests (
    seed_count REAL NOT NULL,
    fulfilled BOOLEAN NOT NULL DEFAULT FALSE,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,

    PRIMARY KEY (user_id, requested_at)
);