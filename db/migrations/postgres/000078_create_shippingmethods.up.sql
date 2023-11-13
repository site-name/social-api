CREATE TABLE IF NOT EXISTS shipping_methods (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(100),
  type character varying(30),
  shipping_zone_id uuid,
  minimum_order_weight real,
  maximum_order_weight real,
  weight_unit character varying(5),
  maximum_delivery_days integer,
  minimum_delivery_days integer,
  description text,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_shipping_methods_name ON shipping_methods USING btree (name);

CREATE INDEX idx_shipping_methods_name_lower_text_pattern ON shipping_methods USING btree (lower((name)::text) text_pattern_ops);