CREATE TABLE IF NOT EXISTS plugin_key_value_store (
  plugin_id character varying(190) NOT NULL PRIMARY KEY,
  p_key character varying(50) NOT NULL,
  p_value bytea,
  expire_at bigint
);