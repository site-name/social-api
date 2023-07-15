ALTER TABLE ONLY warehouseshippingzones
    ADD CONSTRAINT fk_warehouseshippingzones_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES shippingzones(id) ON DELETE CASCADE;
ALTER TABLE ONLY warehouseshippingzones
    ADD CONSTRAINT fk_warehouseshippingzones_warehouses FOREIGN KEY (warehouseid) REFERENCES warehouses(id) ON DELETE CASCADE;
