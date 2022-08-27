CREATE TABLE IF NOT EXISTS salecategories (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  categoryid character varying(36),
  createat bigint
);

ALTER TABLE ONLY salecategories
    ADD CONSTRAINT salecategories_saleid_categoryid_key UNIQUE (saleid, categoryid);

ALTER TABLE ONLY salecategories
    ADD CONSTRAINT fk_salecategories_categories FOREIGN KEY (categoryid) REFERENCES categories(id);

ALTER TABLE ONLY salecategories
    ADD CONSTRAINT fk_salecategories_sales FOREIGN KEY (saleid) REFERENCES sales(id);

