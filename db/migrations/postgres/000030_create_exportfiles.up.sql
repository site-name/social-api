CREATE TABLE IF NOT EXISTS exportfiles (
  id character varying(36) NOT NULL PRIMARY KEY,
  userid character varying(36),
  contentfile text,
  createat bigint,
  updateat bigint
);

ALTER TABLE ONLY exportfiles
    ADD CONSTRAINT fk_exportfiles_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;

