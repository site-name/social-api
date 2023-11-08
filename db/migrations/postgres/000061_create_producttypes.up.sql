CREATE TABLE IF NOT EXISTS product_types (
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

CREATE INDEX idx_product_types_name ON product_types USING btree (name);

CREATE INDEX idx_product_types_name_lower_textpattern ON product_types USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_product_types_slug ON product_types USING btree (slug);
