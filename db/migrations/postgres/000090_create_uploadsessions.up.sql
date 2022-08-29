CREATE TABLE IF NOT EXISTS uploadsessions (
  id character varying(36) NOT NULL PRIMARY KEY,
  type character varying(32),
  createat bigint,
  userid character varying(36),
  filename character varying(256),
  path character varying(512),
  filesize bigint,
  fileoffset bigint
);

CREATE INDEX idx_uploadsessions_create_at ON uploadsessions USING btree (createat);

CREATE INDEX idx_uploadsessions_user_id ON uploadsessions USING btree (type);
