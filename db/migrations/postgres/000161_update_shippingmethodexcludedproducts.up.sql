ALTER TABLE ONLY shipping_method_excluded_products
    ADD CONSTRAINT fk_shipping_method_excluded_products_products FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE ONLY shipping_method_excluded_products
    ADD CONSTRAINT fk_shipping_method_excluded_products_shipping_methods FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods(id);