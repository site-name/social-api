ALTER TABLE ONLY saletranslations
    ADD CONSTRAINT fk_saletranslations_sales FOREIGN KEY (saleid) REFERENCES sales(id) ON DELETE CASCADE;
