CREATE TABLE IF NOT EXISTS wishlists (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  token character varying(36),
  user_id uuid,
  created_at bigint
);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_token_key UNIQUE (token);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_user_id_key UNIQUE (user_id);