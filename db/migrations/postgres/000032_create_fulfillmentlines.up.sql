CREATE TABLE IF NOT EXISTS fulfillment_lines (
  id varchar(36) NOT NULL PRIMARY KEY,
  order_line_id varchar(36) NOT NULL,
  fulfillment_id varchar(36) NOT NULL,
  quantity integer NOT NULL,
  stock_id varchar(36)
);