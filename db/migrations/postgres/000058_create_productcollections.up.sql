CREATE TABLE IF NOT EXISTS productcollections (
  id character varying(36) NOT NULL PRIMARY KEY,
  collectionid character varying(36),
  productid character varying(36)
);

ALTER TABLE ONLY productcollections
    ADD CONSTRAINT productcollections_collectionid_productid_key UNIQUE (collectionid, productid);
