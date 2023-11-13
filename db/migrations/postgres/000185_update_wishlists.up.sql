ALTER TABLE ONLY wishlists
    ADD CONSTRAINT fk_wishlists_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
