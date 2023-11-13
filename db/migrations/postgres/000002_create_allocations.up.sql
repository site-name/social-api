CREATE TABLE IF NOT EXISTS allocations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  order_line_id uuid NOT NULL,
  stock_id uuid NOT NULL,
  quantity_allocated integer NOT NULL
);

ALTER TABLE ONLY allocations
ADD CONSTRAINT allocations_order_line_id_stock_id_key UNIQUE (order_line_id, stock_id);