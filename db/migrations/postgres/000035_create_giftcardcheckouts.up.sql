CREATE TABLE IF NOT EXISTS giftcard_checkouts (
  id character varying(36) NOT NULL PRIMARY KEY,
  giftcardid character varying(36),
  checkoutid character varying(36)
);

