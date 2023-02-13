CREATE TABLE IF NOT EXISTS shippingmethods (
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

CREATE INDEX idx_shipping_methods_name ON shippingmethods USING btree (name);

CREATE INDEX idx_shipping_methods_name_lower_textpattern ON shippingmethods USING btree (lower((name)::text) text_pattern_ops);
ALTER TABLE ONLY shippingmethods
    ADD CONSTRAINT fk_shippingmethods_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES shippingzones(id) ON DELETE CASCADE;
