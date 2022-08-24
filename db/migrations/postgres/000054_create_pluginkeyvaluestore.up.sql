CREATE TABLE IF NOT EXISTS pluginkeyvaluestore (
  pluginid character varying(190) NOT NULL PRIMARY KEY,
  pkey character varying(50) NOT NULL,
  pvalue bytea,
  expireat bigint
);
