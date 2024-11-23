CREATE TABLE IF NOT EXISTS vouchers (
  id varchar(36) NOT NULL PRIMARY KEY,
  type voucher_type NOT NULL,
  name varchar(255),
  usage_limit integer,
  -- used integer NOT NULL,
  start_date bigint NOT NULL, -- begin millisecond of some day
  end_date bigint,
  apply_once_per_order boolean NOT NULL,
  apply_once_per_customer boolean NOT NULL,
  single_use boolean NOT NULL,
  only_for_staff boolean NOT NULL DEFAULT false,
  discount_value_type discount_value_type NOT NULL,
  countries varchar(749) NOT NULL,
  min_checkout_items_quantity integer NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);


-- ALTER TABLE ONLY vouchers
--     ADD CONSTRAINT vouchers_code_key UNIQUE (code);

-- CREATE INDEX idx_vouchers_code ON vouchers USING btree (code);
