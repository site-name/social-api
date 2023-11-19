CREATE TABLE IF NOT EXISTS upload_sessions (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  type character varying(32) NOT NULL,
  created_at bigint NOT NULL,
  user_id uuid NOT NULL,
  file_name character varying(256) NOT NULL,
  path character varying(512) NOT NULL,
  file_size bigint NOT NULL,
  file_offset bigint NOT NULL
);

CREATE INDEX idx_upload_sessions_create_at ON upload_sessions USING btree (created_at);

CREATE INDEX idx_upload_sessions_user_id ON upload_sessions USING btree (user_id);