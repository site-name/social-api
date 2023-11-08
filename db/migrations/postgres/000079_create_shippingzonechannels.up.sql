CREATE TABLE IF NOT EXISTS shipping_zone_channels (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingzoneid character varying(36),
  channelid character varying(36)
);

ALTER TABLE ONLY shipping_zone_channels
    ADD CONSTRAINT shipping_zone_channels_shippingzoneid_channelid_key UNIQUE (shippingzoneid, channelid);

