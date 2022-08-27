CREATE TABLE IF NOT EXISTS roles (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(64),
  displayname character varying(128),
  description character varying(1024),
  createat bigint,
  updateat bigint,
  deleteat bigint,
  permissions text,
  schememanaged boolean,
  builtin boolean
);

ALTER TABLE ONLY roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);
