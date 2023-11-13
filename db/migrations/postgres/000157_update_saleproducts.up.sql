ALTER TABLE ONLY sale_products
    ADD CONSTRAINT fk_sale_products_products FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE ONLY sale_products
    ADD CONSTRAINT fk_sale_products_sales FOREIGN KEY (sale_id) REFERENCES sales(id);
