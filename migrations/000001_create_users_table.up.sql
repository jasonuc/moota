CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS
    users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        username CITEXT NOT NULL UNIQUE,
        email CITEXT NOT NULL UNIQUE,
        title VARCHAR,
        level SMALLINT DEFAULT 1,
        xp INTEGER DEFAULT 0,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );