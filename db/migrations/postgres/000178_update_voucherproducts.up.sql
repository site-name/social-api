ALTER TABLE ONLY voucher_products
    ADD CONSTRAINT fk_voucher_products_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucher_products
    ADD CONSTRAINT fk_voucher_products_vouchers FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;
