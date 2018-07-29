-- +migrate Up
ALTER TABLE ofx_transactions ADD COLUMN trend_id VARCHAR NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE ofx_transactions DROP COLUMN trend_id;
