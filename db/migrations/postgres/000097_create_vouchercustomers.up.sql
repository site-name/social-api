CREATE TABLE IF NOT EXISTS vouchercustomers (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  customeremail character varying(128)
);

ALTER TABLE ONLY vouchercustomers
    ADD CONSTRAINT vouchercustomers_voucherid_customeremail_key UNIQUE (voucherid, customeremail);
