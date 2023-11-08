CREATE TABLE IF NOT EXISTS giftcard_events (
  id character varying(36) NOT NULL PRIMARY KEY,
  date bigint,
  type character varying(255),
  parameters jsonb,
  userid character varying(36),
  giftcardid character varying(36)
);

CREATE INDEX idx_giftcard_events_date ON giftcard_events USING btree (date);

