CREATE TABLE IF NOT EXISTS collectionchannellistings (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  collectionid character varying(36),
  channelid character varying(36),
  publicationdate timestamp with time zone,
  ispublished boolean
);

ALTER TABLE ONLY collectionchannellistings
    ADD CONSTRAINT collectionchannellistings_collectionid_channelid_key UNIQUE (collectionid, channelid);
ALTER TABLE ONLY collectionchannellistings
    ADD CONSTRAINT fk_collectionchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY collectionchannellistings
    ADD CONSTRAINT fk_collectionchannellistings_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;
