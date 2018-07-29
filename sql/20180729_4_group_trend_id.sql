-- +migrate Up
ALTER TABLE transaction_groups ADD COLUMN trend_id VARCHAR NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE transaction_groups DROP COLUMN trend_id;
