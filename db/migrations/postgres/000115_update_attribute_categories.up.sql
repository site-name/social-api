ALTER TABLE ONLY category_attributes
    ADD CONSTRAINT fk_category_attributes_attributes FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY category_attributes
    ADD CONSTRAINT fk_category_attributes_categorys FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;