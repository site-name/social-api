CREATE TABLE IF NOT EXISTS warehouseshippingzones (
  id character varying(36) NOT NULL PRIMARY KEY,
  warehouseid character varying(36),
  shippingzoneid character varying(36)
);

ALTER TABLE ONLY warehouseshippingzones
    ADD CONSTRAINT warehouseshippingzones_warehouseid_shippingzoneid_key UNIQUE (warehouseid, shippingzoneid);
ALTER TABLE ONLY warehouseshippingzones
    ADD CONSTRAINT fk_warehouseshippingzones_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES shippingzones(id);
ALTER TABLE ONLY warehouseshippingzones
    ADD CONSTRAINT fk_warehouseshippingzones_warehouses FOREIGN KEY (warehouseid) REFERENCES warehouses(id);
