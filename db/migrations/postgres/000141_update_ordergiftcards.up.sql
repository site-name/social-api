ALTER TABLE ONLY order_giftcards
    ADD CONSTRAINT fk_order_giftcards_giftcards FOREIGN KEY (giftcard_id) REFERENCES giftcards(id);
ALTER TABLE ONLY order_giftcards
    ADD CONSTRAINT fk_order_giftcards_orders FOREIGN KEY (order_id) REFERENCES orders(id);