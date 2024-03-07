CREATE TABLE IF NOT EXISTS allocations (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  order_line_id varchar(36) NOT NULL,
  stock_id varchar(36) NOT NULL,
  quantity_allocated integer NOT NULL,

  annotations jsonb -- this is a JSON object that can store any additional data you want
);

ALTER TABLE ONLY allocations
ADD CONSTRAINT allocations_order_line_id_stock_id_key UNIQUE (order_line_id, stock_id);