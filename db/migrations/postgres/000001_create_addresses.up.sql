CREATE TABLE IF NOT EXISTS addresses (
  id character varying(36) NOT NULL PRIMARY KEY,
  firstname character varying(64),
  lastname character varying(64),
  companyname character varying(256),
  streetaddress1 character varying(256),
  streetaddress2 character varying(256),
  city character varying(256),
  cityarea character varying(128),
  postalcode character varying(20),
  country character varying(3),
  countryarea character varying(128),
  phone character varying(20),
  createat bigint,
  updateat bigint
);

CREATE INDEX IF NOT EXISTS idx_address_city ON addresses USING btree (city);
CREATE INDEX IF NOT EXISTS idx_address_country ON addresses USING btree (country);
CREATE INDEX IF NOT EXISTS idx_address_firstname ON addresses USING btree (firstname);
CREATE INDEX IF NOT EXISTS idx_address_lastname ON addresses USING btree (lastname);
CREATE INDEX IF NOT EXISTS idx_address_phone ON addresses USING btree (phone);

CREATE INDEX IF NOT EXISTS idx_address_firstname_lower_textpattern ON addresses USING btree (lower((firstname)::text) text_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_address_lastname_lower_textpattern ON addresses USING btree (lower((lastname)::text) text_pattern_ops);
