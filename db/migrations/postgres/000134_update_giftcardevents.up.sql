ALTER TABLE ONLY giftcardevents
    ADD CONSTRAINT fk_giftcardevents_giftcards FOREIGN KEY (giftcardid) REFERENCES giftcards(id) ON DELETE CASCADE;
