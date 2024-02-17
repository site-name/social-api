CREATE TABLE IF NOT EXISTS preferences (
  user_id varchar(36) NOT NULL PRIMARY KEY,
  category varchar(32) NOT NULL,
  name varchar(32) NOT NULL,
  value varchar(2000) NOT NULL
);

CREATE INDEX idx_preferences_category ON preferences USING btree (category);

CREATE INDEX idx_preferences_name ON preferences USING btree (name);

CREATE INDEX idx_preferences_user_id ON preferences USING btree (user_id);