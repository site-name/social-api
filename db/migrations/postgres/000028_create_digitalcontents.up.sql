CREATE TABLE IF NOT EXISTS digitalcontents (
  id character varying(36) NOT NULL PRIMARY KEY,
  shopid character varying(36),
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

ALTER TABLE ONLY digitalcontents
    ADD CONSTRAINT fk_digitalcontents_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;
ALTER TABLE ONLY digitalcontents
    ADD CONSTRAINT fk_digitalcontents_shops FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;
