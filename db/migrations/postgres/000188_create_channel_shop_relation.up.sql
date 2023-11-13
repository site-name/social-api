CREATE TABLE IF NOT EXISTS channel_shops (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  channel_id VARCHAR(36),
  create_at BIGINT,
  end_at BIGINT
);
ALTER TABLE ONLY channel_shops
ADD CONSTRAINT fk_channel_shops_channel FOREIGN KEY (channel_id) REFERENCES channels(id);