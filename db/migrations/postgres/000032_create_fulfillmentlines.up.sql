CREATE TABLE IF NOT EXISTS fulfillment_lines (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  order_line_id uuid NOT NULL,
  fulfillment_id uuid NOT NULL,
  quantity integer NOT NULL,
  stock_id uuid
);