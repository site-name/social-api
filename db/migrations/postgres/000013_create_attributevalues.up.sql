CREATE TABLE IF NOT EXISTS attribute_values (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(250) NOT NULL,
  value varchar(255) NOT NULL,
  slug varchar(255) NOT NULL,
  file_url varchar(200),
  content_type varchar(50),
  attribute_id varchar(36) NOT NULL,
  rich_text text,
  plain_text text,
  "boolean" boolean,
  datetime timestamp with time zone,
  sort_order integer
);

ALTER TABLE ONLY attribute_values
    ADD CONSTRAINT attribute_values_slug_attribute_id_key UNIQUE (slug, attribute_id);

CREATE INDEX idx_attribute_values_name ON attribute_values USING btree (name);

CREATE INDEX idx_attribute_values_name_lower_textpattern ON attribute_values USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_attribute_values_slug ON attribute_values USING btree (slug);
