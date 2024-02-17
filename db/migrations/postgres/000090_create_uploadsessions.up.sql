CREATE TABLE IF NOT EXISTS upload_sessions (
  id varchar(36) NOT NULL PRIMARY KEY,
  type upload_type NOT NULL,
  created_at bigint NOT NULL,
  user_id varchar(36) NOT NULL,
  file_name varchar(256) NOT NULL,
  path varchar(512) NOT NULL,
  file_size bigint NOT NULL,
  file_offset bigint NOT NULL
);

CREATE INDEX idx_upload_sessions_create_at ON upload_sessions USING btree (created_at);

CREATE INDEX idx_upload_sessions_user_id ON upload_sessions USING btree (user_id);