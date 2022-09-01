CREATE TABLE IF NOT EXISTS productvariantchannellistings (
  id character varying(36) NOT NULL PRIMARY KEY,
  variantid character varying(36) NOT NULL,
  channelid character varying(36) NOT NULL,
  currency character varying(3),
  priceamount double precision,
  costpriceamount double precision,
  preorderquantitythreshold integer,
  createat bigint
);

ALTER TABLE ONLY productvariantchannellistings
    ADD CONSTRAINT productvariantchannellistings_variantid_channelid_key UNIQUE (variantid, channelid);

