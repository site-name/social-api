DROP INDEX IF EXISTS idx_order_discounts_name;
DROP INDEX IF EXISTS idx_order_discounts_translated_name;
DROP INDEX IF EXISTS idx_order_discounts_name_lower_textpattern;
DROP INDEX IF EXISTS idx_order_discounts_translated_name_lower_textpattern;
DROP INDEX IF EXISTS idx_order_discounts_voucher_code;

DROP TABLE IF EXISTS order_discounts;