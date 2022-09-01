ALTER TABLE ONLY voucherproducts
    ADD CONSTRAINT fk_voucherproducts_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucherproducts
    ADD CONSTRAINT fk_voucherproducts_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
