CREATE TABLE IF NOT EXISTS voucherchannellistings (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  voucherid character varying(36) NOT NULL,
  channelid character varying(36) NOT NULL,
  discountvalue double precision,
  currency character varying(3),
  minspenamount double precision
);

ALTER TABLE ONLY voucherchannellistings
    ADD CONSTRAINT voucherchannellistings_voucherid_channelid_key UNIQUE (voucherid, channelid);
ALTER TABLE ONLY voucherchannellistings
    ADD CONSTRAINT fk_voucherchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucherchannellistings
    ADD CONSTRAINT fk_voucherchannellistings_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
