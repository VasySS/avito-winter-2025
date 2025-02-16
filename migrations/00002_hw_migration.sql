-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_transfer (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    sender_user_id BIGINT NOT NULL,
    receiver_user_id BIGINT NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (sender_user_id) REFERENCES user_info(id),
    FOREIGN KEY (receiver_user_id) REFERENCES user_info(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_transfer;
-- +goose StatementEnd
