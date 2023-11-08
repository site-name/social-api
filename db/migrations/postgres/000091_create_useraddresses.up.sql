CREATE TABLE IF NOT EXISTS user_addresses (
  id character varying(36) NOT NULL PRIMARY KEY,
  userid character varying(36),
  addressid character varying(36)
);

ALTER TABLE ONLY user_addresses
    ADD CONSTRAINT user_addresses_userid_addressid_key UNIQUE (userid, addressid);

