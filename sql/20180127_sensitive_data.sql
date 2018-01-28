-- +migrate Up
ALTER TABLE ofx_transactions DROP CONSTRAINT ofx_transactions_source_source_account_id_source_account_ty_key,
	DROP COLUMN source,
	DROP COLUMN source_account_id,
	DROP COLUMN source_account_type,
	ADD COLUMN recurring_id VARCHAR,
	ADD CONSTRAINT ofx_transactions_prevent_duplicates UNIQUE (account_id, fi_tid);

ALTER TABLE accounts ADD COLUMN sync BOOLEAN;

-- +migrate Down
ALTER TABLE ofx_transactions ADD CONSTRAINT ofx_transactions_source_source_account_id_source_account_ty_key UNIQUE (source, source_account_id, source_account_type, fi_tid),
	ADD COLUMN source VARCHAR,
	ADD COLUMN source_account_id VARCHAR,
	ADD COLUMN source_account_type VARCHAR,
	DROP CONSTRAINT ofx_transactions_prevent_duplicates;

ALTER TABLE accounts DROP COLUMN sync;
