CREATE TABLE IF NOT EXISTS giftcardcheckouts (
  id character varying(36) NOT NULL PRIMARY KEY,
  giftcardid character varying(36),
  checkoutid character varying(36)
);

