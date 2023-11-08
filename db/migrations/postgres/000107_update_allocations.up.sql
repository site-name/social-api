ALTER TABLE ONLY allocations
ADD CONSTRAINT fk_allocations_order_lines FOREIGN KEY (orderlineid) REFERENCES order_lines(id) ON DELETE CASCADE;
ALTER TABLE ONLY allocations
    ADD CONSTRAINT fk_allocations_stocks FOREIGN KEY (stockid) REFERENCES stocks(id) ON DELETE CASCADE;
