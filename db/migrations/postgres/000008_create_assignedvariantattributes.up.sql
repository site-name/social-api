-- CREATE TABLE IF NOT EXISTS assigned_variant_attributes (
--   id varchar(36) NOT NULL PRIMARY KEY,
--   variant_id varchar(36) NOT NULL,
--   assignment_id varchar(36) NOT NULL
-- );

-- ALTER TABLE ONLY assigned_variant_attributes
--     ADD CONSTRAINT assigned_variant_attributes_variant_id_assignment_id_key UNIQUE (variant_id, assignment_id);