CREATE TABLE IF NOT EXISTS status (
  user_id character varying(36) NOT NULL PRIMARY KEY,
  status character varying(32),
  manual boolean,
  last_activity_at bigint
);

CREATE INDEX idx_status_status ON status USING btree (status);

CREATE INDEX idx_status_user_id ON status USING btree (user_id);