-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_info (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    balance INT NOT NULL CHECK (balance >= 0) DEFAULT 1000,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_info;
-- +goose StatementEnd
