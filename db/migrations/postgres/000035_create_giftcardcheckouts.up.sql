CREATE TABLE IF NOT EXISTS giftcard_checkouts (
  id varchar(36) NOT NULL PRIMARY KEY,
  giftcard_id varchar(36) NOT NULL,
  checkout_id varchar(36) NOT NULL
);