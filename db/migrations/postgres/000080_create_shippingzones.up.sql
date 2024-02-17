CREATE TABLE IF NOT EXISTS shipping_zones (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(100) NOT NULL,
  countries varchar(749) NOT NULL,
  default_flag boolean,
  description varchar(1000) NOT NULL,
  created_at bigint NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_shipping_zones_name ON shipping_zones USING btree (name);

CREATE INDEX idx_shipping_zones_name_lower_text_pattern ON shipping_zones USING btree (lower((name)::text) text_pattern_ops);