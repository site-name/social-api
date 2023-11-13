ALTER TABLE ONLY categories
    ADD CONSTRAINT fk_categories_categories FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE;