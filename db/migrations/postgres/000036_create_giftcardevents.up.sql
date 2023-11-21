CREATE TABLE IF NOT EXISTS giftcard_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  date bigint NOT NULL,
  type GiftcardEventType NOT NULL,
  parameters jsonb,
  user_id uuid,
  giftcard_id uuid NOT NULL
);

CREATE INDEX idx_giftcard_events_date ON giftcard_events USING btree (date);