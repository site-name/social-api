ALTER TABLE ONLY menuitems
    ADD CONSTRAINT fk_menuitems_categories FOREIGN KEY (categoryid) REFERENCES categories(id) ON DELETE CASCADE;
ALTER TABLE ONLY menuitems
    ADD CONSTRAINT fk_menuitems_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY menuitems
    ADD CONSTRAINT fk_menuitems_menuitems FOREIGN KEY (parentid) REFERENCES menuitems(id) ON DELETE CASCADE;
ALTER TABLE ONLY menuitems
    ADD CONSTRAINT fk_menuitems_menus FOREIGN KEY (menuid) REFERENCES menus(id) ON DELETE CASCADE;
ALTER TABLE ONLY menuitems
    ADD CONSTRAINT fk_menuitems_pages FOREIGN KEY (pageid) REFERENCES pages(id) ON DELETE CASCADE;
