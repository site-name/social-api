ALTER TABLE ONLY shipping_methods
    ADD CONSTRAINT fk_shipping_methods_shipping_zones FOREIGN KEY (shipping_zone_id) REFERENCES shipping_zones(id) ON DELETE CASCADE;
