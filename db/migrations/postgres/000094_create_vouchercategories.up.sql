CREATE TABLE IF NOT EXISTS vouchercategories (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  categoryid character varying(36),
  createat bigint
);

ALTER TABLE ONLY vouchercategories
    ADD CONSTRAINT vouchercategories_voucherid_categoryid_key UNIQUE (voucherid, categoryid);
ALTER TABLE ONLY vouchercategories
    ADD CONSTRAINT fk_vouchercategories_categories FOREIGN KEY (categoryid) REFERENCES categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY vouchercategories
    ADD CONSTRAINT fk_vouchercategories_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
