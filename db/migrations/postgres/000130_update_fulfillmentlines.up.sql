ALTER TABLE ONLY fulfillment_lines
    ADD CONSTRAINT fk_fulfillment_lines_order_lines FOREIGN KEY (order_line_id) REFERENCES order_lines(id) ON DELETE CASCADE;
ALTER TABLE ONLY fulfillment_lines
    ADD CONSTRAINT fk_fulfillment_lines_stocks FOREIGN KEY (stock_id) REFERENCES stocks(id);