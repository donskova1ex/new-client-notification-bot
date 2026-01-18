-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages_history (
    message_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    is_sent BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY (message_id, chat_id)
);
CREATE INDEX IF NOT EXISTS idx_messages_history_chat_id ON messages_history (chat_id);
CREATE INDEX IF NOT EXISTS idx_messages_history_created_at ON messages_history (created_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_messages_history_unsent ON messages_history (message_id, chat_id) WHERE is_sent = false AND deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages_history;
-- +goose StatementEnd
