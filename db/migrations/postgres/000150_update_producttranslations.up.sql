ALTER TABLE ONLY producttranslations
    ADD CONSTRAINT fk_producttranslations_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
