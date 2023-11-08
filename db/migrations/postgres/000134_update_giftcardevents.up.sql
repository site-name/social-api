ALTER TABLE ONLY giftcard_events
    ADD CONSTRAINT fk_giftcard_events_giftcards FOREIGN KEY (giftcardid) REFERENCES giftcards(id) ON DELETE CASCADE;
