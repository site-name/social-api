CREATE TABLE IF NOT EXISTS channelshops (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  shopid VARCHAR(36),
  channelid VARCHAR(36),
  createat BIGINT,
  endat BIGINT
);

ALTER TABLE ONLY channelshops
  ADD CONSTRAINT fk_channelshops_shop FOREIGN KEY (shopid) REFERENCES shops(id);

ALTER TABLE ONLY channelshops
  ADD CONSTRAINT fk_channelshops_channel FOREIGN KEY (channelid) REFERENCES channels(id);

ALTER TABLE ONLY channelshops
  ADD CONSTRAINT channelshops_shopid_channelid_key UNIQUE (shopid, channelid);