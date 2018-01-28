-- +migrate Up
CREATE TABLE recurring_groups (
	id 		VARCHAR NOT NULL DEFAULT '',
	name 		VARCHAR NOT NULL DEFAULT '',
	account_ids     VARCHAR[] NOT NULL DEFAULT '{}',
	PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE recurring_groups;
