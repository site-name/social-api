CREATE TABLE IF NOT EXISTS transaction_events (
  id varchar(36) NOT NULL PRIMARY KEY,
  idempotency_key varchar(512),
  psp_reference varchar(512),
  message varchar(512),
  transaction_item_id varchar(36),
  external_url varchar(1000),
  currency Currency NOT NULL,
  type transaction_event_type NOT NULL,
  amount_value decimal(12,3) NOT NULL DEFAULT 0.00,
  user_id varchar(36),
  app_id varchar(36),
  app_identifier varchar(256),
  include_in_calculations boolean NOT NULL DEFAULT false,
  related_granted_refund_id varchar(36),
  created_at bigint NOT NULL
);

ALTER TABLE ONLY transaction_events ADD CONSTRAINT fk_transaction_item_id FOREIGN KEY (transaction_item_id) REFERENCES transaction_items(token) ON DELETE CASCADE;
ALTER TABLE ONLY transaction_events ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
-- ALTER TABLE transaction_events ADD CONSTRAINT fk_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE SET NULL;
CREATE INDEX transaction_events_transaction_item_id_idempotency_key_idx ON transaction_events (transaction_item_id, idempotency_key);
ALTER TABLE ONLY transaction_events ADD CONSTRAINT fk_related_granted_refund_id FOREIGN KEY (related_granted_refund_id) REFERENCES order_granted_refunds(id) ON DELETE SET NULL;
