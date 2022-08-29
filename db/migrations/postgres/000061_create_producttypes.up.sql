CREATE TABLE IF NOT EXISTS producttypes (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  slug character varying(255),
  kind character varying(32),
  hasvariants boolean,
  isshippingrequired boolean,
  isdigital boolean,
  weight real,
  weightunit text,
  metadata jsonb,
  privatemetadata jsonb
);

CREATE INDEX idx_product_types_name ON producttypes USING btree (name);

CREATE INDEX idx_product_types_name_lower_textpattern ON producttypes USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_product_types_slug ON producttypes USING btree (slug);
