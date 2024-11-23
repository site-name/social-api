CREATE TABLE IF NOT EXISTS preorder_reservations (
  id varchar(36) NOT NULL PRIMARY KEY,
  checkout_line_id varchar(36) NOT NULL,
  product_variant_channel_listing_id varchar(36) NOT NULL,
  quantity_reserved integer NOT NULL,
  reserved_until bigint
);

CREATE INDEX IF NOT EXISTS preorder_reservations_checkout_line_id ON preorder_reservations (checkout_line_id);
CREATE INDEX IF NOT EXISTS preorder_reservations_reserved_until ON preorder_reservations (reserved_until);
CREATE UNIQUE INDEx IF NOT EXISTS preorder_reservations_checkout_line_id_product_variant_channel_listing_id ON preorder_reservations (checkout_line_id, product_variant_channel_listing_id);
