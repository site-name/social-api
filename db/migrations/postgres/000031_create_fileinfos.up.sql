CREATE TABLE IF NOT EXISTS file_infos (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  creator_id uuid,
  parent_id uuid,
  created_at bigint,
  updated_at bigint,
  delete_at bigint,
  path character varying(512),
  thumbnail_path character varying(512),
  preview_path character varying(512),
  name character varying(256),
  extension character varying(64),
  size bigint,
  mime_type character varying(256),
  width integer,
  height integer,
  has_preview_image boolean,
  mini_preview bytea,
  content text,
  remote_id character varying(26)
);

CREATE INDEX idx_fileinfo_content_txt ON file_infos USING gin (to_tsvector('english'::regconfig, content));

CREATE INDEX idx_fileinfo_create_at ON file_infos USING btree (create_at);

CREATE INDEX idx_fileinfo_delete_at ON file_infos USING btree (delete_at);

CREATE INDEX idx_fileinfo_extension_at ON file_infos USING btree (extension);

CREATE INDEX idx_fileinfo_name_splitted ON file_infos USING gin (to_tsvector('english'::regconfig, translate((name)::text, '.,-'::text, '   '::text)));

CREATE INDEX idx_fileinfo_name_txt ON file_infos USING gin (to_tsvector('english'::regconfig, (name)::text));

CREATE INDEX idx_fileinfo_parent_id ON file_infos USING btree (parent_id);

CREATE INDEX idx_fileinfo_update_at ON file_infos USING btree (update_at);