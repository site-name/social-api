ALTER TABLE ONLY shippingmethods
    ADD CONSTRAINT fk_shippingmethods_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES shippingzones(id) ON DELETE CASCADE;
