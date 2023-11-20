CREATE TABLE IF NOT EXISTS categories (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(250) NOT NULL,
  slug varchar(255) NOT NULL,
  description jsonb,
  parent_id uuid,
  level smallint NOT NULL,
  background_image varchar(200),
  background_image_alt varchar(128) NOT NULL,
  images varchar(1000),
  seo_title varchar(70) NOT NULL,
  seo_description varchar(300) NOT NULL,
  name_translation jsonb,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);