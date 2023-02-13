CREATE TABLE IF NOT EXISTS sales (
  id character varying(36) NOT NULL PRIMARY KEY,
  shopid character varying(36),
  name character varying(255),
  type character varying(10),
  startdate TIMESTAMP WITH TIME ZONE,
  enddate TIMESTAMP WITH TIME ZONE,
  createat bigint,
  updateat bigint,
  metadata jsonb,
  privatemetadata jsonb
);

CREATE INDEX idx_sales_name ON sales USING btree (name);

CREATE INDEX idx_sales_type ON sales USING btree (type);
ALTER TABLE ONLY sales
    ADD CONSTRAINT fk_sales_shops FOREIGN KEY (shopid) REFERENCES shops(id);
