CREATE TABLE IF NOT EXISTS shippingmethodpostalcoderules (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingmethodid character varying(36),
  start character varying(32),
  'end' character varying(32),
  inclusiontype character varying(32)
);
ALTER TABLE ONLY shippingmethodpostalcoderules
ADD CONSTRAINT shippingmethodpostalcoderules_shippingmethodid_start_end_key UNIQUE (shippingmethodid, 'start', 'end');