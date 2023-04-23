ALTER TABLE ONLY shopstaffs
    ADD CONSTRAINT fk_shopstaffs_users FOREIGN KEY (staffid) REFERENCES users(id);
