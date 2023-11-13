CREATE TABLE IF NOT EXISTS customer_notes (
	id character varying(36) PRIMARY KEY,
	user_id uuid,
	date bigint,
	content text,
	is_public boolean,
	customer_id character varying(36)
);

CREATE INDEX idx_customer_notes_date ON customer_notes USING btree (date);