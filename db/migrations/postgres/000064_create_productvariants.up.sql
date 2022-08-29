CREATE TABLE IF NOT EXISTS productvariants (
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

ALTER TABLE ONLY productvariants
    ADD CONSTRAINT productvariants_sku_key UNIQUE (sku);

CREATE INDEX idx_product_variants_sku ON productvariants USING btree (sku);

ALTER TABLE ONLY productvariants
    ADD CONSTRAINT fk_productvariants_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
