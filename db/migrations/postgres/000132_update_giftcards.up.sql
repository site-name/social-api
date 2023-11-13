ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_products FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_users FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_used_by_users FOREIGN KEY (used_by_id) REFERENCES users(id);