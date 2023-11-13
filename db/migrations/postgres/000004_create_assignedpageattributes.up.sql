CREATE TABLE IF NOT EXISTS assigned_page_attributes (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  page_id uuid NOT NULL,
  assignment_id uuid NOT NULL
);

ALTER TABLE ONLY assigned_page_attributes
    ADD CONSTRAINT assigned_page_attributes_page_id_assignment_id_key UNIQUE (page_id, assignment_id);

