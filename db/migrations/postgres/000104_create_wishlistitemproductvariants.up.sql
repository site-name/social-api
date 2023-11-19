CREATE TABLE IF NOT EXISTS wishlist_item_product_variants (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  wishlist_item_id uuid NOT NULL,
  product_variant_id uuid NOT NULL
);

ALTER TABLE ONLY wishlist_item_product_variants
    ADD CONSTRAINT wishlist_item_product_variants_wishlist_item_id_product_variant_id_key UNIQUE (wishlist_item_id, product_variant_id);
