CREATE TABLE IF NOT EXISTS assigned_product_attributes (
  id varchar(36) NOT NULL PRIMARY KEY,
  product_id varchar(36) NOT NULL,
  assignment_id varchar(36) NOT NULL
);

ALTER TABLE ONLY assigned_product_attributes
    ADD CONSTRAINT assigned_product_attributes_product_id_assignment_id_key UNIQUE (product_id, assignment_id);