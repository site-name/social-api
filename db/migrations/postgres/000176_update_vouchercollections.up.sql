ALTER TABLE ONLY voucher_collections
    ADD CONSTRAINT fk_voucher_collections_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucher_collections
    ADD CONSTRAINT fk_voucher_collections_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
