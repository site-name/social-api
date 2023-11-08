DROP INDEX IF EXISTS idx_order_lines_product_name_lower_textpattern;
DROP INDEX IF EXISTS idx_order_lines_translated_product_name;
DROP INDEX IF EXISTS idx_order_lines_translated_variant_name;
DROP INDEX IF EXISTS idx_order_lines_variant_name;
DROP INDEX IF EXISTS idx_order_lines_variant_name_lower_textpattern;
DROP INDEX IF EXISTS idx_order_lines_product_name;

DROP TABLE IF EXISTS order_lines;