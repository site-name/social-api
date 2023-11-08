CREATE TABLE IF NOT EXISTS shipping_methods (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(100),
  type character varying(30),
  shippingzoneid character varying(36),
  minimumorderweight real,
  maximumorderweight real,
  weightunit character varying(5),
  maximumdeliverydays integer,
  minimumdeliverydays integer,
  description text,
  metadata jsonb,
  privatemetadata jsonb
);

CREATE INDEX idx_shipping_methods_name ON shipping_methods USING btree (name);

CREATE INDEX idx_shipping_methods_name_lower_textpattern ON shipping_methods USING btree (lower((name)::text) text_pattern_ops);
