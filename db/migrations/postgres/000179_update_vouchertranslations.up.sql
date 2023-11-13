ALTER TABLE ONLY voucher_translations
    ADD CONSTRAINT fk_voucher_translations_vouchers FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;
