CREATE TABLE IF NOT EXISTS productchannellistings (
  id character varying(36) NOT NULL PRIMARY KEY,
  productid character varying(36),
  channelid character varying(36),
  visibleinlistings boolean,
  availableforpurchase timestamp with time zone,
  currency character varying(3),
  discountedpriceamount double precision,
  createat bigint,
  publicationdate timestamp with time zone,
  ispublished boolean
);

ALTER TABLE ONLY productchannellistings
    ADD CONSTRAINT productchannellistings_productid_channelid_key UNIQUE (productid, channelid);

CREATE INDEX idx_productchannellistings_puplication_date ON productchannellistings USING btree (publicationdate);

ALTER TABLE ONLY productchannellistings
    ADD CONSTRAINT fk_productchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;

ALTER TABLE ONLY productchannellistings
    ADD CONSTRAINT fk_productchannellistings_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
