CREATE TABLE IF NOT EXISTS wishlist_item_product_variants (
  id character varying(36) NOT NULL PRIMARY KEY,
  wishlistitemid character varying(36),
  productvariantid character varying(36)
);

ALTER TABLE ONLY wishlist_item_product_variants
    ADD CONSTRAINT wishlist_item_product_variants_wishlistitemid_productvariantid_key UNIQUE (wishlistitemid, productvariantid);
