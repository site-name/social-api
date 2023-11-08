ALTER TABLE ONLY voucher_products
    ADD CONSTRAINT fk_voucher_products_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucher_products
    ADD CONSTRAINT fk_voucher_products_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
