CREATE TABLE IF NOT EXISTS preorder_allocations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  order_line_id uuid NOT NULL,
  quantity SMALLINT NOT NULL,
  product_variant_channel_listing_id uuid NOT NULL
);

ALTER TABLE ONLY preorder_allocations
    ADD CONSTRAINT preorder_allocations_order_line_id_product_variant_channel_listing_id_key UNIQUE (order_line_id, product_variant_channel_listing_id);