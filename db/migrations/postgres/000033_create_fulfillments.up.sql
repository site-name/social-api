CREATE TABLE IF NOT EXISTS fulfillments (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  fulfillment_order integer,
  order_id uuid,
  status character varying(32),
  tracking_number character varying(255),
  created_at bigint,
  shipping_refund_amount double precision,
  total_refund_amount double precision,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_fulfillments_status ON fulfillments USING btree (status);

CREATE INDEX idx_fulfillments_tracking_number ON fulfillments USING btree (tracking_number);