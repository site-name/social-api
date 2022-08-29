CREATE TABLE IF NOT EXISTS shippingzones (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(100),
  countries character varying(749),
  "default" boolean,
  description text,
  createat bigint,
  metadata jsonb,
  privatemetadata jsonb
);

CREATE INDEX idx_shipping_zone_name ON shippingzones USING btree (name);

CREATE INDEX idx_shipping_zone_name_lower_textpattern ON shippingzones USING btree (lower((name)::text) text_pattern_ops);

