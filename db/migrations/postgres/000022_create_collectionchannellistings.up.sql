CREATE TABLE IF NOT EXISTS collection_channel_listings (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  collectionid character varying(36),
  channelid character varying(36),
  publicationdate timestamp with time zone,
  ispublished boolean
);

ALTER TABLE ONLY collection_channel_listings
    ADD CONSTRAINT collection_channel_listings_collectionid_channelid_key UNIQUE (collectionid, channelid);
