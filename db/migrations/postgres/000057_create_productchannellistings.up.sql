CREATE TABLE IF NOT EXISTS product_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id uuid NOT NULL,
  channel_id uuid NOT NULL,
  visible_in_listings boolean NOT NULL,
  available_for_purchase bigint, -- future time in milliseconds
  currency varchar(3) NOT NULL,
  discounted_price_amount decimal(12,3),
  created_at bigint NOT NULL,
  publication_date bigint, -- future time in milliseconds
  is_published boolean NOT NULL
);

ALTER TABLE ONLY product_channel_listings
    ADD CONSTRAINT product_channel_listings_product_id_channel_id_key UNIQUE (product_id, channel_id);
