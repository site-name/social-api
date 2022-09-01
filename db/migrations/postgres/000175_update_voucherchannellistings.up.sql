ALTER TABLE ONLY voucherchannellistings
    ADD CONSTRAINT fk_voucherchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucherchannellistings
    ADD CONSTRAINT fk_voucherchannellistings_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
