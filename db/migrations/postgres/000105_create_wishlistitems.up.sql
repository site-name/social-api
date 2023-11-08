CREATE TABLE IF NOT EXISTS wishlist_items (
  id character varying(36) NOT NULL PRIMARY KEY,
  wishlistid character varying(36),
  productid character varying(36),
  createat bigint
);

ALTER TABLE ONLY wishlist_items
    ADD CONSTRAINT wishlist_items_wishlistid_productid_key UNIQUE (wishlistid, productid);

CREATE INDEX idx_wishlist_items ON wishlist_items USING btree (createat);
