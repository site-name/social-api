CREATE TABLE IF NOT EXISTS shipping_zone_channels (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  shipping_zone_id uuid,
  channel_id character varying(36)
);

ALTER TABLE ONLY shipping_zone_channels
    ADD CONSTRAINT shipping_zone_channels_shipping_zone_id_channel_id_key UNIQUE (shipping_zone_id, channel_id);