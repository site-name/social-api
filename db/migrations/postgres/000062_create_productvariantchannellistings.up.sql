CREATE TABLE IF NOT EXISTS product_variant_channel_listings (
  id varchar(36) NOT NULL PRIMARY KEY,
  variant_id varchar(36) NOT NULL,
  channel_id varchar(36) NOT NULL,
  currency Currency,
  price_amount decimal(12,3),
  cost_price_amount decimal(12,3),
  preorder_quantity_threshold integer,
  created_at bigint NOT NULL,

  annotations jsonb
);

ALTER TABLE ONLY product_variant_channel_listings
    ADD CONSTRAINT product_variant_channel_listings_variant_id_channel_id_key UNIQUE (variant_id, channel_id);