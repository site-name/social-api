CREATE TABLE IF NOT EXISTS terms_of_services (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  userid character varying(36),
  text character varying(65535)
);