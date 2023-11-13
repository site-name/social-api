CREATE TABLE IF NOT EXISTS invoice_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  type character varying(255),
  invoice_id uuid,
  order_id uuid,
  user_id uuid,
  parameters text
);