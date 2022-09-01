ALTER TABLE ONLY saleproducts
    ADD CONSTRAINT fk_saleproducts_products FOREIGN KEY (productid) REFERENCES products(id);
ALTER TABLE ONLY saleproducts
    ADD CONSTRAINT fk_saleproducts_sales FOREIGN KEY (saleid) REFERENCES sales(id);
