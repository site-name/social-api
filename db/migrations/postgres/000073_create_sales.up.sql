CREATE TABLE IF NOT EXISTS sales (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(255),
  type character varying(10),
  startdate bigint,
  enddate bigint,
  createat bigint,
  updateat bigint,
  metadata jsonb,
  privatemetadata jsonb
);

CREATE INDEX idx_sales_name ON sales USING btree (name);

CREATE INDEX idx_sales_type ON sales USING btree (type);
