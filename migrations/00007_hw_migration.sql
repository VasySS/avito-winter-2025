-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_user_transfer_sender_user_id ON user_transfer (sender_user_id);
CREATE INDEX idx_user_transfer_receiver_user_id ON user_transfer (receiver_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_user_transfer_sender_user_id;
DROP INDEX idx_user_transfer_receiver_user_id;
-- +goose StatementEnd
