CREATE TABLE IF NOT EXISTS order_giftcards (
  id character varying(36) NOT NULL PRIMARY KEY,
  giftcardid character varying(36),
  orderid character varying(36)
);

ALTER TABLE ONLY order_giftcards
    ADD CONSTRAINT order_giftcards_giftcardid_orderid_key UNIQUE (giftcardid, orderid);
