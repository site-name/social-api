CREATE TABLE IF NOT EXISTS shippingzonechannels (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingzoneid character varying(36),
  channelid character varying(36)
);

ALTER TABLE ONLY shippingzonechannels
    ADD CONSTRAINT shippingzonechannels_shippingzoneid_channelid_key UNIQUE (shippingzoneid, channelid);

ALTER TABLE ONLY shippingzonechannels
    ADD CONSTRAINT fk_shippingzonechannels_channels FOREIGN KEY (channelid) REFERENCES channels(id);

ALTER TABLE ONLY shippingzonechannels
    ADD CONSTRAINT fk_shippingzonechannels_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES shippingzones(id);

