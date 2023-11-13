CREATE TABLE IF NOT EXISTS upload_sessions (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  type character varying(32),
  created_at bigint,
  user_id uuid,
  file_name character varying(256),
  path character varying(512),
  file_size bigint,
  file_offset bigint
);

CREATE INDEX idx_upload_sessions_create_at ON upload_sessions USING btree (create_at);

CREATE INDEX idx_upload_sessions_user_id ON upload_sessions USING btree (user_id);