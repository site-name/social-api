CREATE TABLE IF NOT EXISTS sale_categories (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  categoryid character varying(36),
  createat bigint
);

ALTER TABLE ONLY sale_categories
    ADD CONSTRAINT sale_categories_saleid_categoryid_key UNIQUE (saleid, categoryid);

