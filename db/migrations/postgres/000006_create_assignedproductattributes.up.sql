CREATE TABLE IF NOT EXISTS assignedproductattributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  productid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assignedproductattributes
    ADD CONSTRAINT assignedproductattributes_productid_assignmentid_key UNIQUE (productid, assignmentid);

