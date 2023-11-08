ALTER TABLE ONLY menu_item_translations
    ADD CONSTRAINT fk_menu_item_translations_menu_items FOREIGN KEY (menuitemid) REFERENCES menu_items(id) ON DELETE CASCADE;
