CREATE TABLE IF NOT EXISTS shippingmethodchannellistings (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingmethodid character varying(36),
  channelid character varying(36),
  minimumorderpriceamount double precision,
  currency character varying(3),
  maximumorderpriceamount double precision,
  priceamount double precision,
  createat bigint
);

ALTER TABLE ONLY shippingmethodchannellistings
    ADD CONSTRAINT shippingmethodchannellistings_shippingmethodid_channelid_key UNIQUE (shippingmethodid, channelid);
ALTER TABLE ONLY shippingmethodchannellistings
    ADD CONSTRAINT fk_shippingmethodchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY shippingmethodchannellistings
    ADD CONSTRAINT fk_shippingmethodchannellistings_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES shippingmethods(id) ON DELETE CASCADE;
