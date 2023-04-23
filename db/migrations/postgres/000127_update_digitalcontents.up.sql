ALTER TABLE ONLY digitalcontents
    ADD CONSTRAINT fk_digitalcontents_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;
