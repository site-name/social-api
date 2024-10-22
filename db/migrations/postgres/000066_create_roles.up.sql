CREATE TABLE IF NOT EXISTS roles (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(64) NOT NULL,
  display_name varchar(128) NOT NULL,
  description varchar(1024) NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  delete_at bigint,
  permissions text NOT NULL,
  scheme_managed boolean NOT NULL,
  built_in boolean NOT NULL
);

ALTER TABLE ONLY roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);