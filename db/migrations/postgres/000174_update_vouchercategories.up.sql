ALTER TABLE ONLY vouchercategories
    ADD CONSTRAINT fk_vouchercategories_categories FOREIGN KEY (categoryid) REFERENCES categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY vouchercategories
    ADD CONSTRAINT fk_vouchercategories_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
