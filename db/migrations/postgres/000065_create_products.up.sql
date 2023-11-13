CREATE TABLE IF NOT EXISTS products (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  product_type_id uuid,
  name character varying(250),
  slug character varying(255),
  description text,
  description_plain_text text,
  category_id uuid,
  created_at bigint,
  updated_at bigint,
  charge_taxes boolean,
  weight real,
  weight_unit text,
  default_variant_id uuid,
  rating real,
  metadata jsonb,
  private_metadata jsonb,
  seo_title character varying(70),
  seo_description character varying(300)
);

ALTER TABLE ONLY products
    ADD CONSTRAINT products_name_key UNIQUE (name);

ALTER TABLE ONLY products
    ADD CONSTRAINT products_slug_key UNIQUE (slug);

CREATE INDEX idx_products_metadata ON products USING btree (metadata);

CREATE INDEX idx_products_name ON products USING btree (name);

CREATE INDEX idx_products_name_lower_textpattern ON products USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_products_private_metadata ON products USING btree (private_metadata);

CREATE INDEX idx_products_slug ON products USING btree (slug);