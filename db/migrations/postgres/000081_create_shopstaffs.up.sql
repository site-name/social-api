CREATE TABLE IF NOT EXISTS shop_staffs (
  id character varying(36) NOT NULL PRIMARY KEY,
  staffid character varying(36),
  createat bigint,
  endat bigint,
  salaryperiod character varying(10),
  slary double precision,
  salarycurrency character varying(5)
);
