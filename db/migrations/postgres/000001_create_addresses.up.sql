CREATE TABLE IF NOT EXISTS addresses (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  first_name character varying(64) NOT NULL,
  last_name character varying(64) NOT NULL,
  company_name character varying(256) NOT NULL,
  street_address1 character varying(256) NOT NULL,
  street_address2 character varying(256) NOT NULL,
  city character varying(256) NOT NULL,
  city_area character varying(128) NOT NULL,
  postal_code character varying(20) NOT NULL,
  country character varying(3) NOT NULL,
  country_area character varying(128) NOT NULL,
  phone character varying(11) NOT NULL,
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