CREATE TABLE IF NOT EXISTS wishlist_items (
  id varchar(36) NOT NULL PRIMARY KEY,
  wishlist_id varchar(36) NOT NULL,
  product_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY wishlist_items
    ADD CONSTRAINT wishlist_items_wishlist_id_product_id_key UNIQUE (wishlist_id, product_id);

CREATE INDEX idx_wishlist_items ON wishlist_items USING btree (created_at);
