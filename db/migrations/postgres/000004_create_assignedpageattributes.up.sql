CREATE TABLE IF NOT EXISTS assignedpageattributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  pageid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assignedpageattributes
    ADD CONSTRAINT assignedpageattributes_pageid_assignmentid_key UNIQUE (pageid, assignmentid);

ALTER TABLE ONLY assignedpageattributes
    ADD CONSTRAINT fk_assignedpageattributes_attributepages FOREIGN KEY (assignmentid) REFERENCES attributepages(id) ON DELETE CASCADE;
ALTER TABLE ONLY assignedpageattributes
    ADD CONSTRAINT fk_assignedpageattributes_pages FOREIGN KEY (pageid) REFERENCES pages(id) ON DELETE CASCADE;

