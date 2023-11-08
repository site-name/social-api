CREATE TABLE IF NOT EXISTS assigned_page_attributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  pageid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assigned_page_attributes
    ADD CONSTRAINT assigned_page_attributes_pageid_assignmentid_key UNIQUE (pageid, assignmentid);

