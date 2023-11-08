CREATE TABLE IF NOT EXISTS warehouse_shipping_zones (
  id character varying(36) NOT NULL PRIMARY KEY,
  warehouseid character varying(36),
  shippingzoneid character varying(36)
);

ALTER TABLE ONLY warehouse_shipping_zones
    ADD CONSTRAINT warehouse_shipping_zones_warehouseid_shippingzoneid_key UNIQUE (warehouseid, shippingzoneid);
