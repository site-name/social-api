CREATE TABLE IF NOT EXISTS invoice_events (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  type invoice_event_type NOT NULL,
  invoice_id varchar(36),
  order_id varchar(36),
  user_id varchar(36),
  parameters jsonb,
  app_id varchar(36)
);
