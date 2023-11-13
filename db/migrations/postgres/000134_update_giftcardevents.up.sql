ALTER TABLE ONLY giftcard_events
    ADD CONSTRAINT fk_giftcard_events_giftcards FOREIGN KEY (giftcard_id) REFERENCES giftcards(id) ON DELETE CASCADE;