ALTER TABLE ONLY vouchers
    ADD CONSTRAINT fk_vouchers_shops FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;
