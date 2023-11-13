CREATE TABLE IF NOT EXISTS roles (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(64),
  display_name character varying(128),
  description character varying(1024),
  created_at bigint,
  updated_at bigint,
  delete_at bigint,
  permissions text,
  scheme_managed boolean,
  built_in boolean
);

ALTER TABLE ONLY roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);