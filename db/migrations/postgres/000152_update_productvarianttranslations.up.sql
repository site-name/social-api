ALTER TABLE ONLY productvarianttranslations
    ADD CONSTRAINT fk_productvarianttranslations_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;
