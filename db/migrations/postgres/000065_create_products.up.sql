CREATE TABLE IF NOT EXISTS products (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  product_type_id uuid NOT NULL,
  name character varying(250) NOT NULL,
  slug character varying(255) NOT NULL,
  description jsonb,
  description_plain_text text NOT NULL,
  category_id uuid,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  charge_taxes boolean,
  weight real,
  weight_unit text NOT NULL,
  default_variant_id uuid,
  rating real,
  metadata jsonb,
  private_metadata jsonb,
  seo_title character varying(70) NOT NULL,
  seo_description character varying(300) NOT NULL
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