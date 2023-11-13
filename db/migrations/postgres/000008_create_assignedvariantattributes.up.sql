CREATE TABLE IF NOT EXISTS assigned_variant_attributes (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  variant_id uuid NOT NULL,
  assignment_id uuid NOT NULL
);

ALTER TABLE ONLY assigned_variant_attributes
    ADD CONSTRAINT assigned_variant_attributes_variant_id_assignment_id_key UNIQUE (variant_id, assignment_id);