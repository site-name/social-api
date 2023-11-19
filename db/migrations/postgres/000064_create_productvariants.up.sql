CREATE TABLE IF NOT EXISTS product_variants (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(255) NOT NULL,
  product_id uuid NOT NULL,
  sku character varying(255) NOT NULL,
  weight real,
  weight_unit text NOT NULL,
  track_inventory boolean,
  is_preorder boolean NOT NULL,
  preorder_end_date bigint,
  preorder_global_threshold integer,
  sort_order integer,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY product_variants
    ADD CONSTRAINT product_variants_sku_key UNIQUE (sku);

CREATE INDEX idx_product_variants_sku ON product_variants USING btree (sku);