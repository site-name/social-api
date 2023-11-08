CREATE TABLE IF NOT EXISTS voucher_categories (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  categoryid character varying(36),
  createat bigint
);

ALTER TABLE ONLY voucher_categories
    ADD CONSTRAINT voucher_categories_voucherid_categoryid_key UNIQUE (voucherid, categoryid);
