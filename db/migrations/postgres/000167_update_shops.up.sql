ALTER TABLE ONLY shops
    ADD CONSTRAINT fk_shops_addresses FOREIGN KEY (address_id) REFERENCES addresses(id);
ALTER TABLE ONLY shops
    ADD CONSTRAINT fk_shops_menus FOREIGN KEY (top_menu_id) REFERENCES menus(id);