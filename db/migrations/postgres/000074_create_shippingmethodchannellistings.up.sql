CREATE TABLE IF NOT EXISTS shipping_method_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  shipping_method_id uuid,
  channel_id uuid,
  minimum_order_price_amount double precision,
  currency character varying(3),
  maximum_order_price_amount double precision,
  price_amount double precision,
  created_at bigint
);

ALTER TABLE ONLY shipping_method_channel_listings
    ADD CONSTRAINT shipping_method_channel_listings_shipping_method_id_channel_id_key UNIQUE (shipping_method_id, channel_id);