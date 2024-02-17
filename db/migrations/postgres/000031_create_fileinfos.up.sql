CREATE TABLE IF NOT EXISTS file_infos (
  id varchar(36) NOT NULL PRIMARY KEY,
  creator_id varchar(36) NOT NULL,
  parent_id varchar(36) NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  delete_at bigint,
  path varchar(512) NOT NULL,
  thumbnail_path varchar(512) NOT NULL,
  preview_path varchar(512) NOT NULL,
  name varchar(256) NOT NULL,
  extension varchar(64) NOT NULL,
  size bigint NOT NULL,
  mime_type varchar(256) NOT NULL,
  width integer,
  height integer,
  has_preview_image boolean NOT NULL,
  mini_preview bytea,
  content text NOT NULL,
  remote_id varchar(36)
);

CREATE INDEX idx_fileinfo_content_txt ON file_infos USING gin (to_tsvector('english'::regconfig, content));

CREATE INDEX idx_fileinfo_create_at ON file_infos USING btree (created_at);

CREATE INDEX idx_fileinfo_delete_at ON file_infos USING btree (delete_at);

CREATE INDEX idx_fileinfo_extension_at ON file_infos USING btree (extension);

CREATE INDEX idx_fileinfo_name_splitted ON file_infos USING gin (to_tsvector('english'::regconfig, translate((name)::text, '.,-'::text, '   '::text)));

CREATE INDEX idx_fileinfo_name_txt ON file_infos USING gin (to_tsvector('english'::regconfig, (name)::text));

CREATE INDEX idx_fileinfo_parent_id ON file_infos USING btree (parent_id);

CREATE INDEX idx_fileinfo_update_at ON file_infos USING btree (updated_at);