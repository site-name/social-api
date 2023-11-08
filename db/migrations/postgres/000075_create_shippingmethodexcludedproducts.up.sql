CREATE TABLE IF NOT EXISTS shipping_method_excluded_products (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingmethodid character varying(36),
  productid character varying(36)
);

ALTER TABLE ONLY shipping_method_excluded_products
    ADD CONSTRAINT shipping_method_excluded_products_shippingmethodid_productid_key UNIQUE (shippingmethodid, productid);

