-- +goose Up
ALTER TABLE sessions ADD COLUMN role TEXT NOT NULL DEFAULT 'worker'
    CHECK (role IN ('worker', 'architect', 'orchestrator'));

-- +goose Down
-- SQLite cannot drop a constrained column without rebuilding the table. The
-- role column is additive and older binaries safely ignore it.
