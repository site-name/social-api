CREATE TABLE IF NOT EXISTS warehouse_shipping_zones (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  warehouse_id uuid NOT NULL,
  shipping_zone_id uuid NOT NULL
);

ALTER TABLE ONLY warehouse_shipping_zones
    ADD CONSTRAINT warehouse_shipping_zones_warehouse_id_shipping_zone_id_key UNIQUE (warehouse_id, shipping_zone_id);
