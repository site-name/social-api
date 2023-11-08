CREATE TABLE IF NOT EXISTS file_infos (
  id character varying(36) NOT NULL PRIMARY KEY,
  creatorid character varying(36),
  parentid character varying(36),
  createat bigint,
  updateat bigint,
  deleteat bigint,
  path character varying(512),
  thumbnailpath character varying(512),
  previewpath character varying(512),
  name character varying(256),
  extension character varying(64),
  size bigint,
  mimetype character varying(256),
  width integer,
  height integer,
  haspreviewimage boolean,
  minipreview bytea,
  content text,
  remoteid character varying(26)
);

CREATE INDEX idx_fileinfo_content_txt ON file_infos USING gin (to_tsvector('english'::regconfig, content));

CREATE INDEX idx_fileinfo_create_at ON file_infos USING btree (createat);

CREATE INDEX idx_fileinfo_delete_at ON file_infos USING btree (deleteat);

CREATE INDEX idx_fileinfo_extension_at ON file_infos USING btree (extension);

CREATE INDEX idx_fileinfo_name_splitted ON file_infos USING gin (to_tsvector('english'::regconfig, translate((name)::text, '.,-'::text, '   '::text)));

CREATE INDEX idx_fileinfo_name_txt ON file_infos USING gin (to_tsvector('english'::regconfig, (name)::text));

CREATE INDEX idx_fileinfo_parent_id ON file_infos USING btree (parentid);

CREATE INDEX idx_fileinfo_update_at ON file_infos USING btree (updateat);
