CREATE TABLE IF NOT EXISTS order_events (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  type order_event_type NOT NULL,
  order_id varchar(36) NOT NULL,
  parameters text,
  user_id varchar(36)
);