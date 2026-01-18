-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_types;
-- +goose StatementEnd