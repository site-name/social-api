CREATE TABLE IF NOT EXISTS giftcard_checkouts (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  giftcard_id uuid,
  checkout_id character varying(36)
);