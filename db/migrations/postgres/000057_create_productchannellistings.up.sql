CREATE TABLE IF NOT EXISTS product_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id uuid,
  channel_id uuid,
  visible_in_listings boolean,
  available_for_purchase timestamp with time zone,
  currency character varying(3),
  discounted_price_amount double precision,
  created_at bigint,
  publication_date timestamp with time zone,
  is_published boolean
);

ALTER TABLE ONLY product_channel_listings
    ADD CONSTRAINT product_channel_listings_product_id_channel_id_key UNIQUE (product_id, channel_id);
