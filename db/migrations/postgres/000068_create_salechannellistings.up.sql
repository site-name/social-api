CREATE TABLE IF NOT EXISTS sale_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  sale_id uuid NOT NULL,
  channel_id uuid NOT NULL,
  discount_value decimal(12,3) DEFAULT 0.00,
  currency varchar(3) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_channel_listings
    ADD CONSTRAINT sale_channel_listings_sale_id_channel_id_key UNIQUE (sale_id, channel_id);