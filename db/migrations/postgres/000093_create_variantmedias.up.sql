CREATE TABLE IF NOT EXISTS variant_media (
  id varchar(36) NOT NULL PRIMARY KEY,
  variant_id varchar(36) NOT NULL,
  media_id varchar(36) NOT NULL
);

ALTER TABLE ONLY variant_media
    ADD CONSTRAINT variant_media_variant_id_media_id_key UNIQUE (variant_id, media_id);