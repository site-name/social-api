CREATE TABLE IF NOT EXISTS vouchers (
  id character varying(36) NOT NULL PRIMARY KEY,
  shopid character varying(36),
  type character varying(20),
  name character varying(255),
  code character varying(16),
  usagelimit integer,
  used integer,
  startdate bigint,
  enddate bigint,
  applyonceperorder boolean,
  applyoncepercustomer boolean,
  onlyforstaff boolean,
  discountvaluetype character varying(10),
  countries character varying(749),
  mincheckoutitemsquantity integer,
  createat bigint,
  updateat bigint,
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY vouchers
    ADD CONSTRAINT vouchers_code_key UNIQUE (code);

CREATE INDEX idx_vouchers_code ON vouchers USING btree (code);
ALTER TABLE ONLY vouchers
    ADD CONSTRAINT fk_vouchers_shops FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;
