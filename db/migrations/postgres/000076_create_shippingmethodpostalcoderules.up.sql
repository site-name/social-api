CREATE TABLE IF NOT EXISTS shipping_method_postal_code_rules (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingmethodid character varying(36),
  "start" character varying(32),
  "end" character varying(32),
  inclusiontype character varying(32)
);
ALTER TABLE ONLY shipping_method_postal_code_rules
ADD CONSTRAINT shipping_method_postal_code_rules_shippingmethodid_start_end_key UNIQUE (shippingmethodid, "start", "end");