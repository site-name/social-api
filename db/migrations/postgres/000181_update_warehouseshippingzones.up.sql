ALTER TABLE ONLY warehouseshippingzones
    ADD CONSTRAINT fk_warehouseshippingzones_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES shippingzones(id);
ALTER TABLE ONLY warehouseshippingzones
    ADD CONSTRAINT fk_warehouseshippingzones_warehouses FOREIGN KEY (warehouseid) REFERENCES warehouses(id);
