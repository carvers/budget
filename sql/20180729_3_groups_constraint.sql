-- +migrate Up
ALTER INDEX recurring_groups_pkey RENAME TO transaction_groups_pkey;

-- +migrate Down
ALTER INDEX transaction_groups_pkey RENAME TO recurring_groups_pkey;
