CREATE TABLE IF NOT EXISTS sales (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(255),
  type character varying(10),
  start_date bigint,
  end_date bigint,
  created_at bigint,
  updated_at bigint,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_sales_name ON sales USING btree (name);

CREATE INDEX idx_sales_type ON sales USING btree (type);