ALTER TABLE ONLY sale_categories
    ADD CONSTRAINT fk_sale_categories_categories FOREIGN KEY (category_id) REFERENCES categories(id);
ALTER TABLE ONLY sale_categories
    ADD CONSTRAINT fk_sale_categories_sales FOREIGN KEY (sale_id) REFERENCES sales(id);
