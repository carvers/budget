-- +migrate Up
ALTER TABLE ofx_transactions RENAME recurring_id TO group_id;
ALTER TABLE recurring_groups RENAME TO transaction_groups;

-- +migrate Down
ALTER TABLE ofx_transactions RENAME group_id TO recurring_id;
ALTER TABLE transaction_groups RENAME TO recurring_groups;
