CREATE TABLE IF NOT EXISTS giftcard_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  date bigint,
  type character varying(255),
  parameters jsonb,
  user_id uuid,
  giftcard_id character varying(36)
);

CREATE INDEX idx_giftcard_events_date ON giftcard_events USING btree (date);