CREATE TABLE IF NOT EXISTS order_events (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  type order_event_type NOT NULL,
  order_id varchar(36) NOT NULL,
  parameters jsonb,
  user_id varchar(36),
  app_id varchar(36),
  related_id varchar(36) -- self reference
);

ALTER TABLE order_events ADD CONSTRAINT fk_related_id FOREIGN KEY (related_id) REFERENCES order_events(id) ON DELETE SET NULL;
-- ALTER TABLE order_events ADD CONSTRAINT fk_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE SET NULL;

