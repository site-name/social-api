CREATE TABLE IF NOT EXISTS product_variant_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  variant_id uuid NOT NULL,
  channel_id uuid NOT NULL,
  currency character varying(3),
  price_amount decimal(12,3),
  cost_price_amount decimal(12,3),
  preorder_quantity_threshold integer,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY product_variant_channel_listings
    ADD CONSTRAINT product_variant_channel_listings_variant_id_channel_id_key UNIQUE (variant_id, channel_id);