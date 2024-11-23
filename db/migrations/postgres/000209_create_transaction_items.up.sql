CREATE TABLE IF NOT EXISTS transaction_items (
  token varchar(36) NOT NULL PRIMARY KEY,
  use_old_id boolean NOT NULL DEFAULT false,
  created_at bigint NOT NULL,
  modified_at bigint NOT NULL,
  idempotency_key varchar(512),
  name varchar(512),
  message varchar(512),
  psp_reference varchar(512),
  available_actions varchar(1000),
  currency Currency NOT NULL,
  charged_value decimal(12,3) NOT NULL DEFAULT 0.00,
  authorized_value decimal(12,3) NOT NULL DEFAULT 0.00,
  refunded_value decimal(12,3) NOT NULL DEFAULT 0.00,
  canceled_value decimal(12,3) NOT NULL DEFAULT 0.00,
  refund_pending_value decimal(12,3) NOT NULL DEFAULT 0.00,
  charge_pending_value decimal(12,3) NOT NULL DEFAULT 0.00,
  authorize_pending_value decimal(12,3) NOT NULL DEFAULT 0.00,
  cancel_pending_value decimal(12,3) NOT NULL DEFAULT 0.00,
  external_url varchar(1000),
  checkout_id varchar(36),
  order_id varchar(36),
  user_id varchar(36),
  app_id varchar(36),
  app_identifier varchar(256),
  last_refund_success boolean default true,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY transaction_items ADD CONSTRAINT fk_checkout_id FOREIGN KEY (checkout_id) REFERENCES checkouts(token) ON DELETE SET NULL;
ALTER TABLE ONLY transaction_items ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT;
ALTER TABLE ONLY transaction_items ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
-- ALTER TABLE transaction_items ADD CONSTRAINT fk_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE SET NULL;
CREATE UNIQUE INDEX transaction_items_app_identifier_idempotency_key_idx ON transaction_items (app_identifier, idempotency_key);
