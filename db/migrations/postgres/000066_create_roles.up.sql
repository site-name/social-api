CREATE TABLE IF NOT EXISTS roles (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
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