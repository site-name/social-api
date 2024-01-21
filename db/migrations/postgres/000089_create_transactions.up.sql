CREATE TABLE IF NOT EXISTS transactions (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  payment_id uuid NOT NULL,
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