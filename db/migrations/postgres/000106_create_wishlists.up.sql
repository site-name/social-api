CREATE TABLE IF NOT EXISTS wishlists (
  id character varying(36) NOT NULL PRIMARY KEY,
  token character varying(36),
  userid character varying(36),
  createat bigint
);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_token_key UNIQUE (token);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_userid_key UNIQUE (userid);
