CREATE TABLE IF NOT EXISTS allocations (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  orderlineid character varying(36),
  stockid character varying(36),
  quantityallocated integer
);

ALTER TABLE ONLY allocations
ADD CONSTRAINT allocations_orderlineid_stockid_key UNIQUE (orderlineid, stockid);

ALTER TABLE ONLY allocations
ADD CONSTRAINT fk_allocations_orderlines FOREIGN KEY (orderlineid) REFERENCES orderlines(id) ON DELETE CASCADE;

ALTER TABLE ONLY allocations
    ADD CONSTRAINT fk_allocations_stocks FOREIGN KEY (stockid) REFERENCES stocks(id) ON DELETE CASCADE;
