CREATE TABLE IF NOT EXISTS status (
  userid character varying(36) NOT NULL PRIMARY KEY,
  status character varying(32),
  manual boolean,
  lastactivityat bigint
);

CREATE INDEX idx_status_status ON status USING btree (status);

CREATE INDEX idx_status_user_id ON status USING btree (userid);
