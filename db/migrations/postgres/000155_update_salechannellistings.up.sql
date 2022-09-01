ALTER TABLE ONLY salechannellistings
    ADD CONSTRAINT fk_salechannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY salechannellistings
    ADD CONSTRAINT fk_salechannellistings_sales FOREIGN KEY (saleid) REFERENCES sales(id) ON DELETE CASCADE;
