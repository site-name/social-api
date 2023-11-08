CREATE TABLE IF NOT EXISTS variant_media (
  id character varying(36) NOT NULL PRIMARY KEY,
  variantid character varying(36),
  mediaid character varying(36)
);

ALTER TABLE ONLY variant_media
    ADD CONSTRAINT variant_media_variantid_mediaid_key UNIQUE (variantid, mediaid);
