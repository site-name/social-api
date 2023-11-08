CREATE TABLE IF NOT EXISTS sale_channel_listings (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  channelid character varying(36) NOT NULL,
  discountvalue double precision,
  currency text,
  createat bigint
);

ALTER TABLE ONLY sale_channel_listings
    ADD CONSTRAINT sale_channel_listings_saleid_channelid_key UNIQUE (saleid, channelid);

