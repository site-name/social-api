CREATE TABLE IF NOT EXISTS vouchers (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  type character varying(20) NOT NULL,
  name character varying(255),
  code character varying(16) NOT NULL,
  usage_limit integer,
  used integer NOT NULL,
  start_date bigint NOT NULL, -- begin millisecond of some day
  end_date bigint,
  apply_once_per_order boolean NOT NULL,
  apply_once_per_customer boolean NOT NULL,
  only_for_staff boolean,
  discount_value_type character varying(10) NOT NULL,
  countries character varying(749) NOT NULL,
  min_checkout_items_quantity integer NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY vouchers
    ADD CONSTRAINT vouchers_code_key UNIQUE (code);

CREATE INDEX idx_vouchers_code ON vouchers USING btree (code);
