CREATE TABLE IF NOT EXISTS product_channel_listings (
  id varchar(36) NOT NULL PRIMARY KEY,
  product_id varchar(36) NOT NULL,
  channel_id varchar(36) NOT NULL,
  visible_in_listings boolean NOT NULL,
  available_for_purchase_at bigint, -- future time in milliseconds
  currency Currency NOT NULL,
  discounted_price_amount decimal(12,3),
  created_at bigint NOT NULL,
  publication_date bigint, -- future time in milliseconds
  is_published boolean NOT NULL,
  discounted_price_dirty boolean NOT NULL
);

ALTER TABLE ONLY product_channel_listings
    ADD CONSTRAINT product_channel_listings_product_id_channel_id_key UNIQUE (product_id, channel_id);
