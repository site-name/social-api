CREATE TABLE IF NOT EXISTS shipping_zones (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(100),
  countries character varying(749),
  default_flag boolean,
  description text,
  created_at bigint,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_shipping_zones_name ON shipping_zones USING btree (name);

CREATE INDEX idx_shipping_zones_name_lower_text_pattern ON shipping_zones USING btree (lower((name)::text) text_pattern_ops);