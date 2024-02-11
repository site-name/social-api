-- CREATE TABLE IF NOT EXISTS assigned_variant_attribute_values (
--   id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
--   value_id uuid NOT NULL,
--   assignment_id uuid NOT NULL,
--   sort_order integer
-- );

-- ALTER TABLE ONLY assigned_variant_attribute_values
--     ADD CONSTRAINT assigned_variant_attribute_values_value_id_assignment_id_key UNIQUE (value_id, assignment_id);