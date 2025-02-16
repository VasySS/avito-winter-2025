-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS merch_item (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR NOT NULL UNIQUE,
    price INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS merch_item;
-- +goose StatementEnd
