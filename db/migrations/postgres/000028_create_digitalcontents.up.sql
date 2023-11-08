CREATE TABLE IF NOT EXISTS digital_contents (
  id character varying(36) NOT NULL PRIMARY KEY,
  usedefaultsettings boolean,
  automaticfulfillment boolean,
  contenttype character varying(128),
  productvariantid character varying(36),
  contentfile character varying(200),
  maxdownloads integer,
  urlvaliddays integer,
  metadata jsonb,
  privatemetadata jsonb
);

