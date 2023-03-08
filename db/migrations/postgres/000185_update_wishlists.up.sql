ALTER TABLE ONLY wishlists
    ADD CONSTRAINT fk_wishlists_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;
