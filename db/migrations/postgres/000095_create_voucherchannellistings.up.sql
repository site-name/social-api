CREATE TABLE IF NOT EXISTS voucher_channel_listings (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  voucher_id uuid NOT NULL,
  channel_id uuid NOT NULL,
  discount_value decimal(12,3) NOT NULL DEFAULT 0.00,
  currency Currency NOT NULL,
  min_spend_amount decimal(12,3) NOT NULL DEFAULT 0.00
);

ALTER TABLE ONLY voucher_channel_listings
    ADD CONSTRAINT voucher_channel_listings_voucher_id_channel_id_key UNIQUE (voucher_id, channel_id);