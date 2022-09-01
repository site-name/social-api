ALTER TABLE ONLY fulfillmentlines
    ADD CONSTRAINT fk_fulfillmentlines_orderlines FOREIGN KEY (orderlineid) REFERENCES orderlines(id) ON DELETE CASCADE;
ALTER TABLE ONLY fulfillmentlines
    ADD CONSTRAINT fk_fulfillmentlines_stocks FOREIGN KEY (stockid) REFERENCES stocks(id);
