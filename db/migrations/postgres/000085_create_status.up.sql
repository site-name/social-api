CREATE TABLE IF NOT EXISTS status (
  user_id varchar(36) NOT NULL PRIMARY KEY,
  status varchar(32) NOT NULL,
  manual boolean NOT NULL,
  last_activity_at bigint NOT NULL
);

CREATE INDEX idx_status_status ON status USING btree (status);

CREATE INDEX idx_status_user_id ON status USING btree (user_id);