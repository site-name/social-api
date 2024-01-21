CREATE TABLE IF NOT EXISTS order_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  type order_event_type NOT NULL,
  order_id uuid NOT NULL,
  parameters text,
  user_id uuid
);