CREATE TABLE IF NOT EXISTS categories (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(250) NOT NULL,
  slug character varying(255) NOT NULL,
  description jsonb,
  parent_id uuid,
  level smallint NOT NULL,
  background_image character varying(200),
  background_image_alt character varying(128) NOT NULL,
  images character varying(1000),
  seo_title character varying(70) NOT NULL,
  seo_description character varying(300) NOT NULL,
  name_translation jsonb,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);