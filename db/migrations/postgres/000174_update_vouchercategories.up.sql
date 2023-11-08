ALTER TABLE ONLY voucher_categories
    ADD CONSTRAINT fk_voucher_categories_categories FOREIGN KEY (categoryid) REFERENCES categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucher_categories
    ADD CONSTRAINT fk_voucher_categories_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
