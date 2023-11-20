CREATE TABLE IF NOT EXISTS shipping_method_postal_code_rules (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  shipping_method_id uuid NOT NULL,
  "start" varchar(32) NOT NULL,
  "end" varchar(32) NOT NULL,
  inclusion_type varchar(32) NOT NULL
);
ALTER TABLE ONLY shipping_method_postal_code_rules
ADD CONSTRAINT shipping_method_postal_code_rules_shipping_method_id_start_end_key UNIQUE (shipping_method_id, "start", "end");