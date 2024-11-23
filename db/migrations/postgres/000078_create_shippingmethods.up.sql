CREATE TABLE IF NOT EXISTS shipping_methods (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(100) NOT NULL,
  type shipping_method_type NOT NULL,
  shipping_zone_id varchar(36) NOT NULL,
  minimum_order_weight real,
  maximum_order_weight real,
  weight_unit varchar(5) NOT NULL,
  maximum_delivery_days integer,
  minimum_delivery_days integer,
  description jsonb,
  metadata jsonb,
  private_metadata jsonb,
  tax_class_id varchar(36)
);

CREATE INDEX idx_shipping_methods_name ON shipping_methods USING btree (name);

CREATE INDEX idx_shipping_methods_name_lower_text_pattern ON shipping_methods USING btree (lower((name)::text) text_pattern_ops);