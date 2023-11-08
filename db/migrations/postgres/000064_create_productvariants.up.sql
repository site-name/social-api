CREATE TABLE IF NOT EXISTS product_variants (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(255),
  productid character varying(36),
  sku character varying(255),
  weight real,
  weightunit text,
  trackinventory boolean,
  ispreorder boolean,
  preorderenddate bigint,
  preorderglobalthreshold integer,
  sortorder integer,
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY product_variants
    ADD CONSTRAINT product_variants_sku_key UNIQUE (sku);

CREATE INDEX idx_product_variants_sku ON product_variants USING btree (sku);
