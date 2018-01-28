-- +migrate Up
ALTER TABLE recurring_groups ADD COLUMN finished BOOLEAN NOT NULL DEFAULT false;

-- +migrate Down
ALTER TABLE recurring_groups DROP COLUMN finished;
