CREATE TABLE IF NOT EXISTS wishlists (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  token uuid NOT NULL,
  user_id uuid NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_token_key UNIQUE (token);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_user_id_key UNIQUE (user_id);