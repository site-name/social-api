ALTER TABLE ONLY shippingmethodpostalcoderules
    ADD CONSTRAINT fk_shippingmethodpostalcoderules_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES shippingmethods(id) ON DELETE CASCADE;
