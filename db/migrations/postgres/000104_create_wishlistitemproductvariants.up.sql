CREATE TABLE IF NOT EXISTS wishlist_item_product_variants (
  id varchar(36) NOT NULL PRIMARY KEY,
  wishlist_item_id varchar(36) NOT NULL,
  product_variant_id varchar(36) NOT NULL
);

ALTER TABLE ONLY wishlist_item_product_variants
    ADD CONSTRAINT wishlist_item_product_variants_wishlist_item_id_product_variant_id_key UNIQUE (wishlist_item_id, product_variant_id);
