CREATE TABLE IF NOT EXISTS wishlistitemproductvariants (
  id character varying(36) NOT NULL PRIMARY KEY,
  wishlistitemid character varying(36),
  productvariantid character varying(36)
);

ALTER TABLE ONLY wishlistitemproductvariants
    ADD CONSTRAINT wishlistitemproductvariants_wishlistitemid_productvariantid_key UNIQUE (wishlistitemid, productvariantid);

ALTER TABLE ONLY wishlistitemproductvariants
    ADD CONSTRAINT fk_wishlistitemproductvariants_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;

ALTER TABLE ONLY wishlistitemproductvariants
    ADD CONSTRAINT fk_wishlistitemproductvariants_wishlistitems FOREIGN KEY (wishlistitemid) REFERENCES wishlistitems(id) ON DELETE CASCADE;
