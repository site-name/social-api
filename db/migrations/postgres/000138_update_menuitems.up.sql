ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_categories FOREIGN KEY (categoryid) REFERENCES categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_menu_items FOREIGN KEY (parentid) REFERENCES menu_items(id) ON DELETE CASCADE;
ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_menus FOREIGN KEY (menuid) REFERENCES menus(id) ON DELETE CASCADE;
ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_pages FOREIGN KEY (pageid) REFERENCES pages(id) ON DELETE CASCADE;
