CREATE TABLE IF NOT EXISTS channelshops (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  channelid VARCHAR(36),
  createat BIGINT,
  endat BIGINT
);
ALTER TABLE ONLY channelshops
ADD CONSTRAINT fk_channelshops_channel FOREIGN KEY (channelid) REFERENCES channels(id);