ALTER TABLE ONLY shoptranslations
    ADD CONSTRAINT fk_shoptranslations_shops FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;
