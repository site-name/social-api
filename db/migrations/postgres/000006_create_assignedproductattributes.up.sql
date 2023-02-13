CREATE TABLE IF NOT EXISTS assignedproductattributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  productid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assignedproductattributes
    ADD CONSTRAINT assignedproductattributes_productid_assignmentid_key UNIQUE (productid, assignmentid);

ALTER TABLE ONLY assignedproductattributes
    ADD CONSTRAINT fk_assignedproductattributes_attributeproducts FOREIGN KEY (assignmentid) REFERENCES attributeproducts(id) ON DELETE CASCADE;
ALTER TABLE ONLY assignedproductattributes
    ADD CONSTRAINT fk_assignedproductattributes_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
