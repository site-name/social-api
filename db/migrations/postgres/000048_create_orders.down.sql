DROP INDEX IF EXISTS idx_orders_metadata;
DROP INDEX IF EXISTS idx_orders_private_metadata;
DROP INDEX IF EXISTS idx_orders_user_email;
DROP INDEX IF EXISTS idx_orders_user_email_lower_textpattern;

DROP TABLE IF EXISTS orders;