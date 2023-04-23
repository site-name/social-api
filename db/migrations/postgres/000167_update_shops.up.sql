ALTER TABLE ONLY shops
    ADD CONSTRAINT fk_shops_addresses FOREIGN KEY (addressid) REFERENCES addresses(id);
ALTER TABLE ONLY shops
    ADD CONSTRAINT fk_shops_menus FOREIGN KEY (topmenuid) REFERENCES menus(id);