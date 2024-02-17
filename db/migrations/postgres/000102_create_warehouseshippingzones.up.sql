CREATE TABLE IF NOT EXISTS warehouse_shipping_zones (
  id varchar(36) NOT NULL PRIMARY KEY,
  warehouse_id varchar(36) NOT NULL,
  shipping_zone_id varchar(36) NOT NULL
);

ALTER TABLE ONLY warehouse_shipping_zones
    ADD CONSTRAINT warehouse_shipping_zones_warehouse_id_shipping_zone_id_key UNIQUE (warehouse_id, shipping_zone_id);
