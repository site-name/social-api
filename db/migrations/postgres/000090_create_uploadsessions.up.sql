CREATE TABLE IF NOT EXISTS upload_sessions (
  id character varying(36) NOT NULL PRIMARY KEY,
  type character varying(32),
  createat bigint,
  userid character varying(36),
  filename character varying(256),
  path character varying(512),
  filesize bigint,
  fileoffset bigint
);

CREATE INDEX idx_upload_sessions_create_at ON upload_sessions USING btree (createat);

CREATE INDEX idx_upload_sessions_user_id ON upload_sessions USING btree (type);
