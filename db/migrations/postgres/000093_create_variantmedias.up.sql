CREATE TABLE IF NOT EXISTS variant_media (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  variant_id uuid,
  media_id character varying(36)
);

ALTER TABLE ONLY variant_media
    ADD CONSTRAINT variant_media_variant_id_media_id_key UNIQUE (variant_id, media_id);