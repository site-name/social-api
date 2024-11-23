ALTER TABLE ONLY shipping_methods
  ADD CONSTRAINT fk_shipping_methods_tax_class FOREIGN KEY (tax_class_id) REFERENCES tax_classes(id) ON DELETE SET NULL;
