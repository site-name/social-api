ALTER TABLE ONLY vouchertranslations
    ADD CONSTRAINT fk_vouchertranslations_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
