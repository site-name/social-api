CREATE TABLE IF NOT EXISTS order_giftcards (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  giftcard_id uuid NOT NULL,
  order_id uuid NOT NULL
);

ALTER TABLE ONLY order_giftcards
    ADD CONSTRAINT order_giftcards_giftcard_id_order_id_key UNIQUE (giftcard_id, order_id);