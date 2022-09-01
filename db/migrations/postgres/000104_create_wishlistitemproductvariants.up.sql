CREATE TABLE IF NOT EXISTS wishlistitemproductvariants (
  id character varying(36) NOT NULL PRIMARY KEY,
  wishlistitemid character varying(36),
  productvariantid character varying(36)
);

ALTER TABLE ONLY wishlistitemproductvariants
    ADD CONSTRAINT wishlistitemproductvariants_wishlistitemid_productvariantid_key UNIQUE (wishlistitemid, productvariantid);
