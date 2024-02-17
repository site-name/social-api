CREATE TABLE IF NOT EXISTS customer_events (
  id varchar(36) NOT NULL PRIMARY KEY,
  date bigint NOT NULL,
  type customer_event_type NOT NULL,
  order_id varchar(36),
  user_id varchar(36),
  parameters jsonb
);