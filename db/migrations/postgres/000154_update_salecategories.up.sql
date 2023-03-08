ALTER TABLE ONLY salecategories
    ADD CONSTRAINT fk_salecategories_categories FOREIGN KEY (categoryid) REFERENCES categories(id);
ALTER TABLE ONLY salecategories
    ADD CONSTRAINT fk_salecategories_sales FOREIGN KEY (saleid) REFERENCES sales(id);
