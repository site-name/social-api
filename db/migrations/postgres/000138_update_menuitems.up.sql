ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_categories FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_collections FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_menu_items FOREIGN KEY (parent_id) REFERENCES menu_items(id) ON DELETE CASCADE;
ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_menus FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE;
ALTER TABLE ONLY menu_items
    ADD CONSTRAINT fk_menu_items_pages FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE;