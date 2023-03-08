ALTER TABLE ONLY wishlistitemproductvariants
    ADD CONSTRAINT fk_wishlistitemproductvariants_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;
ALTER TABLE ONLY wishlistitemproductvariants
    ADD CONSTRAINT fk_wishlistitemproductvariants_wishlistitems FOREIGN KEY (wishlistitemid) REFERENCES wishlistitems(id) ON DELETE CASCADE;
