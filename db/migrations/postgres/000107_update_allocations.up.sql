ALTER TABLE ONLY allocations
ADD CONSTRAINT fk_allocations_order_lines FOREIGN KEY (order_line_id) REFERENCES order_lines(id) ON DELETE CASCADE;
ALTER TABLE ONLY allocations
    ADD CONSTRAINT fk_allocations_stocks FOREIGN KEY (stock_id) REFERENCES stocks(id) ON DELETE CASCADE;