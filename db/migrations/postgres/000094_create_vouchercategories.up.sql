CREATE TABLE IF NOT EXISTS vouchercategories (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  categoryid character varying(36),
  createat bigint
);

ALTER TABLE ONLY vouchercategories
    ADD CONSTRAINT vouchercategories_voucherid_categoryid_key UNIQUE (voucherid, categoryid);
