CREATE TABLE IF NOT EXISTS voucher_channel_listings (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  voucherid character varying(36) NOT NULL,
  channelid character varying(36) NOT NULL,
  discountvalue double precision,
  currency character varying(3),
  minspenamount double precision
);

ALTER TABLE ONLY voucher_channel_listings
    ADD CONSTRAINT voucher_channel_listings_voucherid_channelid_key UNIQUE (voucherid, channelid);
