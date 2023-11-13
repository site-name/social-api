CREATE TABLE IF NOT EXISTS voucher_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  voucher_id character varying(36) NOT NULL,
  channel_id character varying(36) NOT NULL,
  discount_value double precision,
  currency character varying(3),
  min_spend_amount double precision
);

ALTER TABLE ONLY voucher_channel_listings
    ADD CONSTRAINT voucher_channel_listings_voucher_id_channel_id_key UNIQUE (voucher_id, channel_id);