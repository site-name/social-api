CREATE TABLE IF NOT EXISTS shipping_method_channel_listings (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingmethodid character varying(36),
  channelid character varying(36),
  minimumorderpriceamount double precision,
  currency character varying(3),
  maximumorderpriceamount double precision,
  priceamount double precision,
  createat bigint
);

ALTER TABLE ONLY shipping_method_channel_listings
    ADD CONSTRAINT shipping_method_channel_listings_shippingmethodid_channelid_key UNIQUE (shippingmethodid, channelid);
