CREATE TABLE IF NOT EXISTS customer_notes (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid,
	date bigint NOT NULL,
	content text,
	is_public boolean,
	customer_id uuid NOT NULL
);

CREATE INDEX idx_customer_notes_date ON customer_notes USING btree (date);