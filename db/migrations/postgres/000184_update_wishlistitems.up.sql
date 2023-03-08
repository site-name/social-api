ALTER TABLE ONLY wishlistitems
    ADD CONSTRAINT fk_wishlistitems_productvariants FOREIGN KEY (productid) REFERENCES productvariants(id) ON DELETE CASCADE;
ALTER TABLE ONLY wishlistitems
    ADD CONSTRAINT fk_wishlistitems_wishlists FOREIGN KEY (wishlistid) REFERENCES wishlists(id) ON DELETE CASCADE;
