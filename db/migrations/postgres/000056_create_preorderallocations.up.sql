CREATE TABLE IF NOT EXISTS preorderallocations (
  id character varying(36) NOT NULL PRIMARY KEY,
  orderlineid character varying(36),
  quantity integer,
  productvariantchannellistingid character varying(36)
);

ALTER TABLE ONLY public.preorderallocations
    ADD CONSTRAINT preorderallocations_orderlineid_productvariantchannellistin_key UNIQUE (orderlineid, productvariantchannellistingid);
