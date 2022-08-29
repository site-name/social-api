CREATE TABLE IF NOT EXISTS shopstaffs (
  id character varying(36) NOT NULL PRIMARY KEY,
  shopid character varying(36),
  staffid character varying(36),
  createat bigint,
  endat bigint
);

ALTER TABLE ONLY shopstaffs
    ADD CONSTRAINT shopstaffs_shopid_staffid_key UNIQUE (shopid, staffid);

ALTER TABLE ONLY shopstaffs
    ADD CONSTRAINT fk_shopstaffs_shops FOREIGN KEY (shopid) REFERENCES shops(id);

ALTER TABLE ONLY shopstaffs
    ADD CONSTRAINT fk_shopstaffs_users FOREIGN KEY (staffid) REFERENCES users(id);

