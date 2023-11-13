CREATE TABLE IF NOT EXISTS vouchers (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  type character varying(20),
  name character varying(255),
  code character varying(16),
  usage_limit integer,
  used integer,
  start_date bigint,
  end_date bigint,
  apply_once_per_order boolean,
  apply_once_per_customer boolean,
  only_for_staff boolean,
  discount_value_type character varying(10),
  countries character varying(749),
  min_checkout_items_quantity integer,
  created_at bigint,
  updated_at bigint,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY vouchers
    ADD CONSTRAINT vouchers_code_key UNIQUE (code);

CREATE INDEX idx_vouchers_code ON vouchers USING btree (code);
