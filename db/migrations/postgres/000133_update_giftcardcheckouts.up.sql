ALTER TABLE ONLY giftcard_checkouts
    ADD CONSTRAINT fk_giftcard_checkouts_checkouts FOREIGN KEY (checkout_id) REFERENCES checkouts(token);
ALTER TABLE ONLY giftcard_checkouts
    ADD CONSTRAINT fk_giftcard_checkouts_giftcards FOREIGN KEY (giftcard_id) REFERENCES giftcards(id);