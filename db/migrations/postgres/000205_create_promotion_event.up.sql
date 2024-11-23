CREATE TABLE IF NOT EXISTS promotion_events (
  id varchar(36) NOT NULL PRIMARY KEY,
  date bigint NOT NULL,
  type promotion_event_type NOT NULL,
  parameters jsonb,
  user_id varchar(36),
  app_id varchar(36),
  promotion_id varchar(36)
);

ALTER TABLE promotion_events ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
-- ALTER TABLE promotion_events ADD CONSTRAINT fk_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE SET NULL;
ALTER TABLE promotion_events ADD CONSTRAINT fk_promotion_id FOREIGN KEY (promotion_id) REFERENCES promotions(id) ON DELETE CASCADE;
