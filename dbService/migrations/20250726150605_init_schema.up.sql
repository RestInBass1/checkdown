-- migrations/000001_init_schema.up.sql

CREATE TABLE IF NOT EXISTS tasks (
     id BIGSERIAL PRIMARY KEY,
     title TEXT NOT NULL,
     description TEXT,
     is_done TEXT NOT NULL,
     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
