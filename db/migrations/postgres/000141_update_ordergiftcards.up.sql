ALTER TABLE ONLY order_giftcards
    ADD CONSTRAINT fk_order_giftcards_giftcards FOREIGN KEY (giftcardid) REFERENCES giftcards(id);
ALTER TABLE ONLY order_giftcards
    ADD CONSTRAINT fk_order_giftcards_orders FOREIGN KEY (orderid) REFERENCES orders(id);
