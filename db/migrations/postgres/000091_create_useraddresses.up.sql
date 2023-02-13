CREATE TABLE IF NOT EXISTS useraddresses (
  id character varying(36) NOT NULL PRIMARY KEY,
  userid character varying(36),
  addressid character varying(36)
);

ALTER TABLE ONLY useraddresses
    ADD CONSTRAINT useraddresses_userid_addressid_key UNIQUE (userid, addressid);

ALTER TABLE ONLY useraddresses
    ADD CONSTRAINT fk_useraddresses_addresses FOREIGN KEY (addressid) REFERENCES addresses(id) ON DELETE CASCADE;
ALTER TABLE ONLY useraddresses
    ADD CONSTRAINT fk_useraddresses_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;
