CREATE TABLE IF NOT EXISTS shipping_method_postal_code_rules (
  id varchar(36) NOT NULL PRIMARY KEY,
  shipping_method_id varchar(36) NOT NULL,
  "start" varchar(32) NOT NULL,
  "end" varchar(32) NOT NULL,
  inclusion_type inclusion_type NOT NULL
);
ALTER TABLE ONLY shipping_method_postal_code_rules
ADD CONSTRAINT shipping_method_postal_code_rules_shipping_method_id_start_end_key UNIQUE (shipping_method_id, "start", "end");