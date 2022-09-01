ALTER TABLE ONLY orderlines
    ADD CONSTRAINT fk_orderlines_orders FOREIGN KEY (orderid) REFERENCES orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY orderlines
    ADD CONSTRAINT fk_orderlines_productvariants FOREIGN KEY (variantid) REFERENCES productvariants(id);
