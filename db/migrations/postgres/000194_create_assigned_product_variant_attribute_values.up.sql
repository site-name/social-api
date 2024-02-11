CREATE TABLE IF NOT EXISTS assigned_product_variant_attribute_values (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    attribute_value_id uuid NOT NULL,
    variant_id uuid NOT NULL
);

ALTER TABLE assigned_product_variant_attribute_values
ADD CONSTRAINT assigned_product_variant_attribute_values_attribute_value_id_fk
FOREIGN KEY (attribute_value_id)
REFERENCES custom_product_attribute_values(id) ON DELETE CASCADE;

ALTER TABLE assigned_product_variant_attribute_values
ADD CONSTRAINT assigned_product_variant_attribute_values_variant_id_fk
FOREIGN KEY (variant_id)
REFERENCES product_variants(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX assigned_product_variant_attribute_values_attribute_value_id_variant_id_unique_key
ON assigned_product_variant_attribute_values (attribute_value_id, variant_id);
