ALTER TABLE ONLY warehouse_shipping_zones
    ADD CONSTRAINT fk_warehouse_shipping_zones_shipping_zones FOREIGN KEY (shipping_zone_id) REFERENCES shipping_zones(id) ON DELETE CASCADE;
ALTER TABLE ONLY warehouse_shipping_zones
    ADD CONSTRAINT fk_warehouse_shipping_zones_warehouses FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE;
