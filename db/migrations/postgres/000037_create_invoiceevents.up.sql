CREATE TABLE IF NOT EXISTS invoice_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  type varchar(255) NOT NULL,
  invoice_id uuid,
  order_id uuid,
  user_id uuid,
  parameters jsonb
);