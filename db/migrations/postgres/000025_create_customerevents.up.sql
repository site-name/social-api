CREATE TABLE IF NOT EXISTS customer_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  date bigint NOT NULL,
  type CustomerEventType NOT NULL,
  order_id uuid,
  user_id uuid,
  parameters jsonb
);