ALTER TABLE ONLY sales
    ADD CONSTRAINT fk_sales_shops FOREIGN KEY (shopid) REFERENCES shops(id);
