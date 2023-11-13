CREATE TABLE IF NOT EXISTS sale_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  sale_id uuid,
  channel_id character varying(36) NOT NULL,
  discount_value double precision,
  currency text,
  created_at bigint
);

ALTER TABLE ONLY sale_channel_listings
    ADD CONSTRAINT sale_channel_listings_sale_id_channel_id_key UNIQUE (sale_id, channel_id);