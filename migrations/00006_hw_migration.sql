-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_merch_purchase_item_id ON merch_purchase (merch_item_id);
CREATE INDEX idx_merch_purchase_user_id ON merch_purchase (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_merch_purchase_item_id;
DROP INDEX idx_merch_purchase_user_id;
-- +goose StatementEnd
