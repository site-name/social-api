CREATE TABLE IF NOT EXISTS category_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code varchar(10) NOT NULL,
  category_id uuid NOT NULL,
  name varchar(250) NOT NULL,
  description text NOT NULL,
  seo_title varchar(70),
  seo_description varchar(300)
);

CREATE INDEX idx_category_translations_name ON category_translations USING btree (name);

CREATE INDEX idx_category_translations_name_lower_textpattern ON category_translations USING btree (lower((name)::text) text_pattern_ops);