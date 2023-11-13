CREATE TABLE IF NOT EXISTS fulfillment_lines (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  order_line_id uuid,
  fulfillment_id uuid,
  quantity integer,
  stock_id character varying(36)
);