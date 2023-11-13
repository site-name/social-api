CREATE TABLE IF NOT EXISTS assigned_product_attributes (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id uuid NOT NULL,
  assignment_id uuid NOT NULL
);

ALTER TABLE ONLY assigned_product_attributes
    ADD CONSTRAINT assigned_product_attributes_product_id_assignment_id_key UNIQUE (product_id, assignment_id);