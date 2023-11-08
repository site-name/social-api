CREATE TABLE IF NOT EXISTS assigned_product_attributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  productid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assigned_product_attributes
    ADD CONSTRAINT assigned_product_attributes_productid_assignmentid_key UNIQUE (productid, assignmentid);

