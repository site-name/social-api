ALTER TABLE ONLY productmedias
    ADD CONSTRAINT fk_productmedias_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
