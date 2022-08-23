CREATE TABLE IF NOT EXISTS attributepages (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  pagetypeid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY attributepages
    ADD CONSTRAINT attributepages_attributeid_pagetypeid_key UNIQUE (attributeid, pagetypeid);

ALTER TABLE ONLY attributepages
    ADD 
ALTER TABLE ONLY attributepages
    ADD CONSTRAINT fk_attributepages_pagetypes FOREIGN KEY (pagetypeid) REFERENCES pagetypes(id) ON DELETE CASCADE;

