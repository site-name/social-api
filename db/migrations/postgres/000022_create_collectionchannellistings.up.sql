CREATE TABLE IF NOT EXISTS collection_channel_listings (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  collection_id varchar(36) NOT NULL,
  channel_id varchar(36),
  publication_date timestamp with time zone,
  is_published boolean
);

ALTER TABLE ONLY collection_channel_listings
    ADD CONSTRAINT collection_channel_listings_collection_id_channel_id_key UNIQUE (collection_id, channel_id);