ALTER TABLE ONLY menuitemtranslations
    ADD CONSTRAINT fk_menuitemtranslations_menuitems FOREIGN KEY (menuitemid) REFERENCES menuitems(id) ON DELETE CASCADE;
