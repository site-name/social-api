ALTER TABLE ONLY giftcard_checkouts
    ADD CONSTRAINT fk_giftcard_checkouts_checkouts FOREIGN KEY (checkoutid) REFERENCES checkouts(token);
ALTER TABLE ONLY giftcard_checkouts
    ADD CONSTRAINT fk_giftcard_checkouts_giftcards FOREIGN KEY (giftcardid) REFERENCES giftcards(id);
