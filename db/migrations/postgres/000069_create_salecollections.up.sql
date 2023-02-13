CREATE TABLE IF NOT EXISTS salecollections (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  collectionid character varying(36),
  createat bigint
);

ALTER TABLE ONLY salecollections
    ADD CONSTRAINT salecollections_saleid_collectionid_key UNIQUE (saleid, collectionid);

ALTER TABLE ONLY salecollections
    ADD CONSTRAINT fk_salecollections_collections FOREIGN KEY (collectionid) REFERENCES collections(id);
ALTER TABLE ONLY salecollections
    ADD CONSTRAINT fk_salecollections_sales FOREIGN KEY (saleid) REFERENCES sales(id);
