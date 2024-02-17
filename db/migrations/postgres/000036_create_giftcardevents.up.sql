CREATE TABLE IF NOT EXISTS giftcard_events (
  id varchar(36) NOT NULL PRIMARY KEY,
  date bigint NOT NULL,
  type giftcard_event_type NOT NULL,
  parameters jsonb,
  user_id varchar(36),
  giftcard_id varchar(36) NOT NULL
);

CREATE INDEX idx_giftcard_events_date ON giftcard_events USING btree (date);