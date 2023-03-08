CREATE TABLE IF NOT EXISTS useraddresses (
  id character varying(36) NOT NULL PRIMARY KEY,
  userid character varying(36),
  addressid character varying(36)
);

ALTER TABLE ONLY useraddresses
    ADD CONSTRAINT useraddresses_userid_addressid_key UNIQUE (userid, addressid);

