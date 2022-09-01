CREATE TABLE IF NOT EXISTS salechannellistings (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  channelid character varying(36) NOT NULL,
  discountvalue double precision,
  currency text,
  createat bigint
);

ALTER TABLE ONLY salechannellistings
    ADD CONSTRAINT salechannellistings_saleid_channelid_key UNIQUE (saleid, channelid);

