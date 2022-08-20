DROP INDEX IF EXISTS idx_address_city;
DROP INDEX IF EXISTS idx_address_country ;
DROP INDEX IF EXISTS idx_address_firstname;
DROP INDEX IF EXISTS idx_address_lastname;
DROP INDEX IF EXISTS idx_address_phone;

DROP INDEX IF EXISTS idx_address_firstname_lower_textpattern;
DROP INDEX IF EXISTS idx_address_lastname_lower_textpattern;

DROP TABLE IF EXISTS addresses;