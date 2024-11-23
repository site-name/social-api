CREATE TABLE IF NOT EXISTS reservations (
  id varchar(36) NOT NULL PRIMARY KEY,
  checkout_line_id varchar(36) NOT NULL,
  stock_id varchar(36) NOT NULL,
  quantity_reserved integer NOT NULL,
  reserved_until bigint
);

CREATE INDEX IF NOT EXISTS reservations_checkout_line_id ON reservations (checkout_line_id);
CREATE INDEX IF NOT EXISTS reservations_reserved_until ON reservations (reserved_until);
CREATE UNIQUE INDEx IF NOT EXISTS reservations_checkout_line_id_stock_id ON reservations (checkout_line_id, stock_id);
