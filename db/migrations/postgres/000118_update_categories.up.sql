ALTER TABLE ONLY categories
    ADD CONSTRAINT fk_categories_categories FOREIGN KEY (parentid) REFERENCES categories(id) ON DELETE CASCADE;
