ALTER TABLE ONLY giftcardcheckouts
    ADD CONSTRAINT fk_giftcardcheckouts_checkouts FOREIGN KEY (checkoutid) REFERENCES checkouts(token);
ALTER TABLE ONLY giftcardcheckouts
    ADD CONSTRAINT fk_giftcardcheckouts_giftcards FOREIGN KEY (giftcardid) REFERENCES giftcards(id);
