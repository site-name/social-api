ALTER TABLE ONLY wishlist_items
    ADD CONSTRAINT fk_wishlist_items_product_variants FOREIGN KEY (product_id) REFERENCES product_variants(id) ON DELETE CASCADE;
ALTER TABLE ONLY wishlist_items
    ADD CONSTRAINT fk_wishlist_items_wishlists FOREIGN KEY (wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE;
