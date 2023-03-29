ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_products FOREIGN KEY (productid) REFERENCES products(id);
ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_users FOREIGN KEY (createdbyid) REFERENCES users(id);
ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_used_by_users FOREIGN KEY (usedbyid) REFERENCES users(id);
ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_shopid FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;
