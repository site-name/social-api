ALTER TABLE ONLY sale_categories
    ADD CONSTRAINT fk_sale_categories_categories FOREIGN KEY (categoryid) REFERENCES categories(id);
ALTER TABLE ONLY sale_categories
    ADD CONSTRAINT fk_sale_categories_sales FOREIGN KEY (saleid) REFERENCES sales(id);
