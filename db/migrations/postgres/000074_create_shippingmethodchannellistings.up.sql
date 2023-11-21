CREATE TABLE IF NOT EXISTS shipping_method_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  shipping_method_id uuid NOT NULL,
  channel_id uuid NOT NULL,
  minimum_order_price_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  currency Currency NOT NULL,
  maximum_order_price_amount decimal(12,3),
  price_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY shipping_method_channel_listings
    ADD CONSTRAINT shipping_method_channel_listings_shipping_method_id_channel_id_key UNIQUE (shipping_method_id, channel_id);