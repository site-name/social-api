ALTER TABLE ONLY shipping_methods
    ADD CONSTRAINT fk_shipping_methods_shipping_zones FOREIGN KEY (shippingzoneid) REFERENCES shipping_zones(id) ON DELETE CASCADE;
