CREATE TABLE IF NOT EXISTS export_files (
  id character varying(36) NOT NULL PRIMARY KEY,
  userid character varying(36),
  contentfile text,
  createat bigint,
  updateat bigint
);

