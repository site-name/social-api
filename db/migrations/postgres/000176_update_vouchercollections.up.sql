ALTER TABLE ONLY voucher_collections
    ADD CONSTRAINT fk_voucher_collections_collections FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucher_collections
    ADD CONSTRAINT fk_voucher_collections_vouchers FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;
