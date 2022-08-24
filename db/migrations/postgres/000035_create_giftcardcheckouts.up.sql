CREATE TABLE IF NOT EXISTS giftcardcheckouts (
  id character varying(36) NOT NULL PRIMARY KEY,
  giftcardid character varying(36),
  checkoutid character varying(36)
);

ALTER TABLE ONLY giftcardcheckouts
    ADD CONSTRAINT fk_giftcardcheckouts_checkouts FOREIGN KEY (checkoutid) REFERENCES checkouts(token);

ALTER TABLE ONLY giftcardcheckouts
    ADD CONSTRAINT fk_giftcardcheckouts_giftcards FOREIGN KEY (giftcardid) REFERENCES giftcards(id);

