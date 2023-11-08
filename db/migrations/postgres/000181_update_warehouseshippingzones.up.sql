ALTER TABLE ONLY warehouse_shipping_zones
    ADD CONSTRAINT fk_warehouse_shipping_zones_shipping_zones FOREIGN KEY (shippingzoneid) REFERENCES shipping_zones(id) ON DELETE CASCADE;
ALTER TABLE ONLY warehouse_shipping_zones
    ADD CONSTRAINT fk_warehouse_shipping_zones_warehouses FOREIGN KEY (warehouseid) REFERENCES warehouses(id) ON DELETE CASCADE;
