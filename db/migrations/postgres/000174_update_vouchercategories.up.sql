ALTER TABLE ONLY voucher_categories
    ADD CONSTRAINT fk_voucher_categories_categories FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucher_categories
    ADD CONSTRAINT fk_voucher_categories_vouchers FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;
