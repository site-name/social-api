CREATE TABLE IF NOT EXISTS stocks (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  warehouseid character varying(36),
  productvariantid character varying(36),
  quantity integer
);

ALTER TABLE ONLY stocks
    ADD CONSTRAINT stocks_warehouseid_productvariantid_key UNIQUE (warehouseid, productvariantid);

ALTER TABLE ONLY stocks
    ADD CONSTRAINT fk_stocks_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;

ALTER TABLE ONLY stocks
    ADD CONSTRAINT fk_stocks_warehouses FOREIGN KEY (warehouseid) REFERENCES warehouses(id) ON DELETE CASCADE;
