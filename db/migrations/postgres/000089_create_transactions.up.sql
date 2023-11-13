CREATE TABLE IF NOT EXISTS transactions (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  payment_id uuid,
  token character varying(512),
  kind character varying(25),
  is_success boolean,
  action_required boolean,
  action_required_data text,
  currency character varying(3),
  amount double precision,
  error character varying(256),
  customer_id character varying(256),
  gateway_response text,
  already_processed boolean
);