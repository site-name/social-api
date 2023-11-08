ALTER TABLE ONLY wishlist_item_product_variants
    ADD CONSTRAINT fk_wishlist_item_product_variants_product_variants FOREIGN KEY (productvariantid) REFERENCES product_variants(id) ON DELETE CASCADE;
ALTER TABLE ONLY wishlist_item_product_variants
    ADD CONSTRAINT fk_wishlist_item_product_variants_wishlist_items FOREIGN KEY (wishlistitemid) REFERENCES wishlist_items(id) ON DELETE CASCADE;
