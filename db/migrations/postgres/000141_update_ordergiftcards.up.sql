ALTER TABLE ONLY ordergiftcards
    ADD CONSTRAINT fk_ordergiftcards_giftcards FOREIGN KEY (giftcardid) REFERENCES giftcards(id);
ALTER TABLE ONLY ordergiftcards
    ADD CONSTRAINT fk_ordergiftcards_orders FOREIGN KEY (orderid) REFERENCES orders(id);
