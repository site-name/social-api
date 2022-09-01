CREATE TABLE IF NOT EXISTS categories (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  slug character varying(255),
  description text,
  parentid character varying(36),
  backgroundimage character varying(200),
  backgroundimagealt character varying(128),
  seotitle character varying(70),
  seodescription character varying(300),
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);

