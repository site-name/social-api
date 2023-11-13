CREATE TABLE IF NOT EXISTS product_types (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(250),
  slug character varying(255),
  kind character varying(32),
  has_variants boolean,
  is_shipping_required boolean,
  is_digital boolean,
  weight real,
  weight_unit text,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_product_types_name ON product_types USING btree (name);

CREATE INDEX idx_product_types_name_lower_textpattern ON product_types USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_product_types_slug ON product_types USING btree (slug);