-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS merch_purchase (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    merch_item_id BIGINT NOT NULL,
    price INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user_info (id),
    FOREIGN KEY (merch_item_id) REFERENCES merch_item (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS merch_purchase;
-- +goose StatementEnd
