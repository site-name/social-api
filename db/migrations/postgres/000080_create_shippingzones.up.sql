CREATE TABLE IF NOT EXISTS shipping_zones (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(100),
  countries character varying(749),
  "default" boolean,
  description text,
  createat bigint,
  metadata jsonb,
  privatemetadata jsonb
);

CREATE INDEX idx_shipping_zone_name ON shipping_zones USING btree (name);

CREATE INDEX idx_shipping_zone_name_lower_textpattern ON shipping_zones USING btree (lower((name)::text) text_pattern_ops);

