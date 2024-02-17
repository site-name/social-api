CREATE TABLE IF NOT EXISTS addresses (
  id varchar(36) NOT NULL PRIMARY KEY,
  first_name varchar(64) NOT NULL,
  last_name varchar(64) NOT NULL,
  company_name varchar(256) NOT NULL,
  street_address1 varchar(256) NOT NULL,
  street_address2 varchar(256) NOT NULL,
  city varchar(256) NOT NULL,
  city_area varchar(128) NOT NULL,
  postal_code varchar(20) NOT NULL,
  country country_code NOT NULL, -- enum
  country_area varchar(128) NOT NULL,
  phone varchar(11) NOT NULL,
  user_id varchar(36) NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_address_city ON addresses USING btree (city);
CREATE INDEX IF NOT EXISTS idx_address_country ON addresses USING btree (country);
CREATE INDEX IF NOT EXISTS idx_address_firstname ON addresses USING btree (first_name);
CREATE INDEX IF NOT EXISTS idx_address_lastname ON addresses USING btree (last_name);
CREATE INDEX IF NOT EXISTS idx_address_phone ON addresses USING btree (phone);

CREATE INDEX IF NOT EXISTS idx_address_firstname_lower_textpattern ON addresses USING btree (lower((first_name)::text) text_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_address_lastname_lower_textpattern ON addresses USING btree (lower((last_name)::text) text_pattern_ops);