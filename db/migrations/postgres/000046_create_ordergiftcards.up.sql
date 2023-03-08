CREATE TABLE IF NOT EXISTS ordergiftcards (
  id character varying(36) NOT NULL PRIMARY KEY,
  giftcardid character varying(36),
  orderid character varying(36)
);

ALTER TABLE ONLY ordergiftcards
    ADD CONSTRAINT ordergiftcards_giftcardid_orderid_key UNIQUE (giftcardid, orderid);
