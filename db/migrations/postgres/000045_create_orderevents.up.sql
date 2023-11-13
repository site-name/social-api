CREATE TABLE IF NOT EXISTS order_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  type character varying(255),
  order_id uuid,
  parameters text,
  user_id character varying(36)
);