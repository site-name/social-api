CREATE TABLE IF NOT EXISTS warehouse_shipping_zones (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  warehouse_id uuid,
  shipping_zone_id character varying(36)
);

ALTER TABLE ONLY warehouse_shipping_zones
    ADD CONSTRAINT warehouse_shipping_zones_warehouse_id_shipping_zone_id_key UNIQUE (warehouse_id, shipping_zone_id);
