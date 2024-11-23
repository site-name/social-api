CREATE TABLE IF NOT EXISTS order_granted_refunds (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  amount_value decimal(12,3) NOT NULL DEFAULT 0.00,
  currency Currency NOT NULL,
  reason text,
  user_id varchar(36),
  app_id varchar(36),
  order_id varchar(36),
  shipping_costs_included boolean NOT NULL DEFAULT false,
  transaction_item_id varchar(36),
  status order_granted_refund_status
);

ALTER TABLE order_granted_refunds ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
ALTER TABLE order_granted_refunds ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
-- ALTER TABLE order_granted_refunds ADD CONSTRAINT fk_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_order_granted_refunds_updated_at ON order_granted_refunds USING btree (updated_at);

CREATE TABLE IF NOT EXISTS order_granted_refund_lines (
  id varchar(36) NOT NULL PRIMARY KEY,
  order_line_id varchar(36),
  quantity integer NOT NULL,
  granted_refund_id varchar(36),
  reason text
);

ALTER TABLE order_granted_refund_lines ADD CONSTRAINT fk_granted_refund_id FOREIGN KEY (granted_refund_id) REFERENCES order_granted_refunds(id) ON DELETE CASCADE;
ALTER TABLE order_granted_refund_lines ADD CONSTRAINT fk_order_line_id FOREIGN KEY (order_line_id) REFERENCES order_lines(id) ON DELETE CASCADE;
