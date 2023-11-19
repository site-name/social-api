CREATE TABLE IF NOT EXISTS preferences (
  user_id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  category character varying(32) NOT NULL,
  name character varying(32) NOT NULL,
  value character varying(2000)
);

CREATE INDEX idx_preferences_category ON preferences USING btree (category);

CREATE INDEX idx_preferences_name ON preferences USING btree (name);

CREATE INDEX idx_preferences_user_id ON preferences USING btree (user_id);