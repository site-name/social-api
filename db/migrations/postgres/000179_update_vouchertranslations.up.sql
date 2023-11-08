ALTER TABLE ONLY voucher_translations
    ADD CONSTRAINT fk_voucher_translations_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
