CREATE TABLE IF NOT EXISTS export_files (
  id varchar(36) NOT NULL PRIMARY KEY,
  user_id varchar(36),
  content_file text,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL
);