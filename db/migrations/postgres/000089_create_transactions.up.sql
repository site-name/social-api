CREATE TABLE IF NOT EXISTS payment_transactions (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  payment_id varchar(36) NOT NULL,
  token varchar(512) NOT NULL,
  kind transaction_kind NOT NULL,
  is_success boolean NOT NULL,
  action_required boolean NOT NULL,
  action_required_data text NOT NULL,
  currency Currency NOT NULL,
  amount decimal(12,3) NOT NULL DEFAULT 0.00,
  error varchar(256),
  customer_id varchar(256),
  gateway_response text NOT NULL,
  already_processed boolean NOT NULL
);

ALTER TABLE ONLY payment_transactions
    ADD CONSTRAINT fk_transactions_payments FOREIGN KEY (payment_id) REFERENCES payments(id);
