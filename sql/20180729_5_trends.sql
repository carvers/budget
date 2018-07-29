-- +migrate Up
CREATE TABLE trends (
	id 		   VARCHAR NOT NULL DEFAULT '',
	name 		   VARCHAR NOT NULL DEFAULT '',
	forced_period_days INT     NOT NULL DEFAULT 0,
	PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE trends;
