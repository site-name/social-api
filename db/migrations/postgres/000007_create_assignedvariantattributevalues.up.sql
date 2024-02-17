-- CREATE TABLE IF NOT EXISTS assigned_variant_attribute_values (
--   id varchar(36) NOT NULL PRIMARY KEY,
--   value_id varchar(36) NOT NULL,
--   assignment_id varchar(36) NOT NULL,
--   sort_order integer
-- );

-- ALTER TABLE ONLY assigned_variant_attribute_values
--     ADD CONSTRAINT assigned_variant_attribute_values_value_id_assignment_id_key UNIQUE (value_id, assignment_id);