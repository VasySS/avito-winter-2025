-- +goose Up
-- +goose StatementBegin
INSERT INTO merch_item (name, price, created_at)
VALUES
  ('t-shirt', 80, NOW()),
  ('cup', 20, NOW()),
  ('book', 50, NOW()),
  ('pen', 10, NOW()),
  ('powerbank', 200, NOW()),
  ('hoody', 300, NOW()),
  ('umbrella', 200, NOW()),
  ('socks', 10, NOW()),
  ('wallet', 50, NOW()),
  ('pink-hoody', 500, NOW())
ON CONFLICT (name) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM merch_item;
-- +goose StatementEnd
