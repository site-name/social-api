CREATE TABLE IF NOT EXISTS productmedias (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  productid character varying(36),
  ppoi character varying(20),
  image character varying(200),
  alt character varying(128),
  type character varying(32),
  externalurl character varying(256),
  oembeddata text,
  sortorder integer
);
