CREATE TABLE IF NOT EXISTS categories (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  slug character varying(255),
  description jsonb,
  parentid character varying(36),
  level smallint,
  backgroundimage character varying(200),
  backgroundimagealt character varying(128),
  images character varying(1000),
  seotitle character varying(70),
  seodescription character varying(300),
  nametranslation jsonb,
  metadata jsonb,
  privatemetadata jsonb
);
-- ALTER TABLE ONLY categories
-- ADD CONSTRAINT categories_name_key UNIQUE (name);
ALTER TABLE ONLY categories
ADD CONSTRAINT categories_slug_key UNIQUE (slug);