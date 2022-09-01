CREATE TABLE IF NOT EXISTS warehouseshippingzones (
  id character varying(36) NOT NULL PRIMARY KEY,
  warehouseid character varying(36),
  shippingzoneid character varying(36)
);

ALTER TABLE ONLY warehouseshippingzones
    ADD CONSTRAINT warehouseshippingzones_warehouseid_shippingzoneid_key UNIQUE (warehouseid, shippingzoneid);
