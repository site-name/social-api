CREATE TABLE IF NOT EXISTS assignedpageattributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  pageid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assignedpageattributes
    ADD CONSTRAINT assignedpageattributes_pageid_assignmentid_key UNIQUE (pageid, assignmentid);

