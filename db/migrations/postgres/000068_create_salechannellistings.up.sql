CREATE TABLE IF NOT EXISTS sale_channel_listings (
  id varchar(36) NOT NULL PRIMARY KEY,
  sale_id varchar(36) NOT NULL,
  channel_id varchar(36) NOT NULL,
  discount_value decimal(12,3) DEFAULT 0.00,
  currency Currency NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_channel_listings
    ADD CONSTRAINT sale_channel_listings_sale_id_channel_id_key UNIQUE (sale_id, channel_id);