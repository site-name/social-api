CREATE TABLE IF NOT EXISTS transactions (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  payment_id uuid NOT NULL,
  token character varying(512) NOT NULL,
  kind character varying(25) NOT NULL,
  is_success boolean NOT NULL,
  action_required boolean NOT NULL,
  action_required_data text NOT NULL,
  currency character varying(3) NOT NULL,
  amount double precision,
  error character varying(256),
  customer_id character varying(256),
  gateway_response text NOT NULL,
  already_processed boolean NOT NULL
);