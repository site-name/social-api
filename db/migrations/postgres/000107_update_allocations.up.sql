ALTER TABLE ONLY allocations
ADD CONSTRAINT fk_allocations_orderlines FOREIGN KEY (orderlineid) REFERENCES orderlines(id) ON DELETE CASCADE;
ALTER TABLE ONLY allocations
    ADD CONSTRAINT fk_allocations_stocks FOREIGN KEY (stockid) REFERENCES stocks(id) ON DELETE CASCADE;
