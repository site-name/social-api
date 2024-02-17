CREATE TABLE IF NOT EXISTS wishlists (
  id varchar(36) NOT NULL PRIMARY KEY,
  token varchar(36) NOT NULL,
  user_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_token_key UNIQUE (token);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_user_id_key UNIQUE (user_id);