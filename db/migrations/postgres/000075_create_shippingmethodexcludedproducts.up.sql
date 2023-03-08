CREATE TABLE IF NOT EXISTS shippingmethodexcludedproducts (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingmethodid character varying(36),
  productid character varying(36)
);

ALTER TABLE ONLY shippingmethodexcludedproducts
    ADD CONSTRAINT shippingmethodexcludedproducts_shippingmethodid_productid_key UNIQUE (shippingmethodid, productid);

