CREATE TABLE IF NOT EXISTS wishlistitems (
  id character varying(36) NOT NULL PRIMARY KEY,
  wishlistid character varying(36),
  productid character varying(36),
  createat bigint
);

ALTER TABLE ONLY wishlistitems
    ADD CONSTRAINT wishlistitems_wishlistid_productid_key UNIQUE (wishlistid, productid);

CREATE INDEX idx_wishlist_items ON wishlistitems USING btree (createat);
