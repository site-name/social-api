ALTER TABLE ONLY fulfillment_lines
    ADD CONSTRAINT fk_fulfillment_lines_order_lines FOREIGN KEY (orderlineid) REFERENCES order_lines(id) ON DELETE CASCADE;
ALTER TABLE ONLY fulfillment_lines
    ADD CONSTRAINT fk_fulfillment_lines_stocks FOREIGN KEY (stockid) REFERENCES stocks(id);
