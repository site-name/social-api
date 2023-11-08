ALTER TABLE ONLY sale_translations
    ADD CONSTRAINT fk_sale_translations_sales FOREIGN KEY (saleid) REFERENCES sales(id) ON DELETE CASCADE;
