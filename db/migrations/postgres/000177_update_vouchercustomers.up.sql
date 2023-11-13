ALTER TABLE ONLY voucher_customers
    ADD CONSTRAINT fk_voucher_customers_vouchers FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;
