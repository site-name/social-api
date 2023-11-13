ALTER TABLE ONLY shop_staffs
    ADD CONSTRAINT fk_shop_staffs_users FOREIGN KEY (staff_id) REFERENCES users(id);
