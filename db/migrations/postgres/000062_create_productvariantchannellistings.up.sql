CREATE TABLE IF NOT EXISTS product_variant_channel_listings (
  id character varying(36) NOT NULL PRIMARY KEY,
  variantid character varying(36) NOT NULL,
  channelid character varying(36) NOT NULL,
  currency character varying(3),
  priceamount double precision,
  costpriceamount double precision,
  preorderquantitythreshold integer,
  createat bigint
);

ALTER TABLE ONLY product_variant_channel_listings
    ADD CONSTRAINT product_variant_channel_listings_variantid_channelid_key UNIQUE (variantid, channelid);

