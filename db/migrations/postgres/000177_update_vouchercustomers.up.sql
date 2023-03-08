ALTER TABLE ONLY vouchercustomers
    ADD CONSTRAINT fk_vouchercustomers_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
