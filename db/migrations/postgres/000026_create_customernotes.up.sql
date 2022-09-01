CREATE TABLE IF NOT EXISTS customernotes (
  id character varying(36) PRIMARY KEY,
	userid character varying(36),
	date bigint,
	content text,
	ispublic boolean,
	customerid character varying(36)
);

CREATE INDEX idx_customer_notes_date ON customernotes USING btree (date);
		