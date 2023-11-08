ALTER TABLE ONLY voucher_customers
    ADD CONSTRAINT fk_voucher_customers_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
