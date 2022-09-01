ALTER TABLE ONLY variantmedias
    ADD CONSTRAINT fk_variantmedias_productmedias FOREIGN KEY (mediaid) REFERENCES productmedias(id) ON DELETE CASCADE;
ALTER TABLE ONLY variantmedias
    ADD CONSTRAINT fk_variantmedias_productvariants FOREIGN KEY (variantid) REFERENCES productvariants(id) ON DELETE CASCADE;
