-- +migrate Up
CREATE TABLE accounts (
	id 		VARCHAR,
	name 		VARCHAR,
	account_type 	VARCHAR,
	request_type	VARCHAR,
	bank_org	VARCHAR,
	bank_fid	VARCHAR,
	bank_url	VARCHAR,
	PRIMARY KEY(id)
);

ALTER TABLE ofx_transactions DROP CONSTRAINT ofx_transactions_pkey,
	ADD PRIMARY KEY (id),
	ADD COLUMN account_id VARCHAR;

-- +migrate Down
DROP TABLE accounts;
ALTER TABLE ofx_transactions DROP CONSTRAINT ofx_transactions_pkey,
	ADD PRIMARY KEY (id, source, source_account_id, source_account_type),
	DROP COLUMN account_id;
