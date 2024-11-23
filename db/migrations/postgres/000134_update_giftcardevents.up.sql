ALTER TABLE ONLY giftcard_events
    ADD CONSTRAINT fk_giftcard_events_giftcards FOREIGN KEY (giftcard_id) REFERENCES giftcards(id) ON DELETE CASCADE;

ALTER TABLE ONLY giftcard_events
    ADD CONSTRAINT fk_giftcard_events_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;

ALTER TABLE ONLY giftcard_events
    ADD CONSTRAINT fk_giftcard_events_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
