CREATE TABLE IF NOT EXISTS custom_product_attributes (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(250) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    product_id uuid NOT NULL
);

ALTER TABLE custom_product_attributes
ADD CONSTRAINT custom_product_attributes_product_id_fk 
FOREIGN KEY (product_id) 
REFERENCES products (id) ON DELETE CASCADE;
