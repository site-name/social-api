CREATE TABLE IF NOT EXISTS sales (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(255) NOT NULL,
  type discount_value_type NOT NULL,
  start_date bigint NOT NULL, -- future time in milli
  end_date bigint,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_sales_name ON sales USING btree (name);

CREATE INDEX idx_sales_type ON sales USING btree (type);