CREATE TABLE IF NOT EXISTS saleproductvariants (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  productvariantid character varying(36),
  createat bigint
);

ALTER TABLE ONLY public.saleproductvariants
    ADD CONSTRAINT saleproductvariants_saleid_productvariantid_key UNIQUE (saleid, productvariantid);
