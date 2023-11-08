CREATE TABLE IF NOT EXISTS attribute_pages (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  pagetypeid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY attribute_pages
    ADD CONSTRAINT attribute_pages_attributeid_pagetypeid_key UNIQUE (attributeid, pagetypeid);
