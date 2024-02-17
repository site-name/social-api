CREATE TABLE IF NOT EXISTS shipping_zone_channels (
  id varchar(36) NOT NULL PRIMARY KEY,
  shipping_zone_id varchar(36) NOT NULL,
  channel_id varchar(36) NOT NULL
);

ALTER TABLE ONLY shipping_zone_channels
    ADD CONSTRAINT shipping_zone_channels_shipping_zone_id_channel_id_key UNIQUE (shipping_zone_id, channel_id);