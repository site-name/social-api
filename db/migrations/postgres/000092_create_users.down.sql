DROP INDEX IF EXISTS idx_users_all_no_full_name_txt;
DROP INDEX IF EXISTS idx_users_all_txt;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_email_lower_textpattern;
DROP INDEX IF EXISTS idx_users_firstname_lower_textpattern;
DROP INDEX IF EXISTS idx_users_lastname_lower_textpattern;
DROP INDEX IF EXISTS idx_users_metadata;
DROP INDEX IF EXISTS idx_users_names_no_full_name_txt;
DROP INDEX IF EXISTS idx_users_names_txt;
DROP INDEX IF EXISTS idx_users_nickname_lower_textpattern;
DROP INDEX IF EXISTS idx_users_private_metadata;
DROP INDEX IF EXISTS idx_users_username_lower_textpattern;
DROP INDEX IF EXISTS order_user_search_gin;

DROP TABLE IF EXISTS users;