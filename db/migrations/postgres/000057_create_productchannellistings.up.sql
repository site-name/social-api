CREATE TABLE IF NOT EXISTS product_channel_listings (
  id character varying(36) NOT NULL PRIMARY KEY,
  productid character varying(36),
  channelid character varying(36),
  visibleinlistings boolean,
  availableforpurchase timestamp with time zone,
  currency character varying(3),
  discountedpriceamount double precision,
  createat bigint,
  publicationdate timestamp with time zone,
  ispublished boolean
);

ALTER TABLE ONLY product_channel_listings
    ADD CONSTRAINT product_channel_listings_productid_channelid_key UNIQUE (productid, channelid);

CREATE INDEX idx_product_channel_listings_puplication_date ON product_channel_listings USING btree (publicationdate);
