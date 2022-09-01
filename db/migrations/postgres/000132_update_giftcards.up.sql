ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_products FOREIGN KEY (productid) REFERENCES products(id);
ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_users FOREIGN KEY (createdbyid) REFERENCES users(id);
