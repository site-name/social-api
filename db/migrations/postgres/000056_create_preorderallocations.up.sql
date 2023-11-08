CREATE TABLE IF NOT EXISTS preorder_allocations (
  id character varying(36) NOT NULL PRIMARY KEY,
  orderlineid character varying(36),
  quantity integer,
  productvariantchannellistingid character varying(36)
);

ALTER TABLE ONLY preorder_allocations
    ADD CONSTRAINT preorder_allocations_orderlineid_productvariantchannellistin_key UNIQUE (orderlineid, productvariantchannellistingid);
