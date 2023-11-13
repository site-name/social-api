ALTER TABLE ONLY order_lines
    ADD CONSTRAINT fk_order_lines_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY order_lines
    ADD CONSTRAINT fk_order_lines_product_variants FOREIGN KEY (variant_id) REFERENCES product_variants(id);