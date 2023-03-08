ALTER TABLE ONLY digitalcontents
    ADD CONSTRAINT fk_digitalcontents_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;
ALTER TABLE ONLY digitalcontents
    ADD CONSTRAINT fk_digitalcontents_shops FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;
