ALTER TABLE ONLY assignedpageattributes
    ADD CONSTRAINT fk_assignedpageattributes_attributepages FOREIGN KEY (assignmentid) REFERENCES attributepages(id) ON DELETE CASCADE;
ALTER TABLE ONLY assignedpageattributes
    ADD CONSTRAINT fk_assignedpageattributes_pages FOREIGN KEY (pageid) REFERENCES pages(id) ON DELETE CASCADE;
