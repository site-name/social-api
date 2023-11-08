ALTER TABLE ONLY sale_products
    ADD CONSTRAINT fk_sale_products_products FOREIGN KEY (productid) REFERENCES products(id);
ALTER TABLE ONLY sale_products
    ADD CONSTRAINT fk_sale_products_sales FOREIGN KEY (saleid) REFERENCES sales(id);
