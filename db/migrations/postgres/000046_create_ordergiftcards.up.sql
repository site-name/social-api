CREATE TABLE IF NOT EXISTS order_giftcards (
  id varchar(36) NOT NULL PRIMARY KEY,
  giftcard_id varchar(36) NOT NULL,
  order_id varchar(36) NOT NULL
);

ALTER TABLE ONLY order_giftcards
    ADD CONSTRAINT order_giftcards_giftcard_id_order_id_key UNIQUE (giftcard_id, order_id);