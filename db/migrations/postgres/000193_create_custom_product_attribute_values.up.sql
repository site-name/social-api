CREATE TABLE IF NOT EXISTS custom_product_attribute_values (
    id varchar(36) NOT NULL PRIMARY KEY,
    value VARCHAR(250) NOT NULL,
    attribute_id varchar(36) NOT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (
       SELECT conname
       FROM pg_constraint
       WHERE conname = 'custom_product_attribute_values_attribute_id_fk'
    )
    THEN
        ALTER TABLE custom_product_attribute_values
        ADD CONSTRAINT custom_product_attribute_values_attribute_id_fk
        FOREIGN KEY (attribute_id)
        REFERENCES custom_product_attributes (id) ON DELETE CASCADE;
    END IF;
END $$;
