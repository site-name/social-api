ALTER TABLE ONLY vouchercollections
    ADD CONSTRAINT fk_vouchercollections_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY vouchercollections
    ADD CONSTRAINT fk_vouchercollections_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
