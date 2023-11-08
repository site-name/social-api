CREATE TABLE IF NOT EXISTS attribute_values (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  value character varying(9),
  slug character varying(255),
  fileurl character varying(200),
  contenttype character varying(50),
  attributeid character varying(36),
  richtext text,
  "boolean" boolean,
  datetime timestamp with time zone,
  sortorder integer
);

ALTER TABLE ONLY attribute_values
    ADD CONSTRAINT attribute_values_slug_attributeid_key UNIQUE (slug, attributeid);

CREATE INDEX idx_attribute_values_name ON attribute_values USING btree (name);

CREATE INDEX idx_attribute_values_name_lower_textpattern ON attribute_values USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_attribute_values_slug ON attribute_values USING btree (slug);

