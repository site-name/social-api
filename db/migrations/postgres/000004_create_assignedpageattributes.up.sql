CREATE TABLE IF NOT EXISTS assigned_page_attributes (
  id varchar(36) NOT NULL PRIMARY KEY,
  page_id varchar(36) NOT NULL,
  assignment_id varchar(36) NOT NULL
);

ALTER TABLE ONLY assigned_page_attributes
    ADD CONSTRAINT assigned_page_attributes_page_id_assignment_id_key UNIQUE (page_id, assignment_id);

