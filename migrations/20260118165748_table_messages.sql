-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    planned_date TIMESTAMP,
    user_id BIGINT NOT NULL,
    type_id INT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_planned_date ON messages (planned_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_id ON messages (user_id) WHERE deleted_at IS NULL; 
CREATE INDEX IF NOT EXISTS idx_type_id ON messages (type_id) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
