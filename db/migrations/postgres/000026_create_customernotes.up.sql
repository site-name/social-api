CREATE TABLE IF NOT EXISTS customer_notes (
	id varchar(36) NOT NULL PRIMARY KEY,
	user_id varchar(36),
	date bigint NOT NULL,
	content text,
	is_public boolean,
	customer_id varchar(36) NOT NULL
);

CREATE INDEX idx_customer_notes_date ON customer_notes USING btree (date);