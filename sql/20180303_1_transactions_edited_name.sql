-- +migrate Up
ALTER TABLE ofx_transactions ADD COLUMN edited_name VARCHAR NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE ofx_transactions DROP COLUMN edited_name;
