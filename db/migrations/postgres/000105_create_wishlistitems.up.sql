CREATE TABLE IF NOT EXISTS wishlist_items (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  wishlist_id uuid,
  product_id uuid,
  created_at bigint
);

ALTER TABLE ONLY wishlist_items
    ADD CONSTRAINT wishlist_items_wishlist_id_product_id_key UNIQUE (wishlist_id, product_id);

CREATE INDEX idx_wishlist_items ON wishlist_items USING btree (create_at);
