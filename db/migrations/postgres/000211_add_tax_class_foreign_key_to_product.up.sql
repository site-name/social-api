ALTER TABLE products ADD CONSTRAINT fk_tax_class_id FOREIGN KEY (tax_class_id) REFERENCES tax_classes(id) ON DELETE SET NULL;
