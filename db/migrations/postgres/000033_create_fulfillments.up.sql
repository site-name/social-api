CREATE TABLE IF NOT EXISTS fulfillments (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  fulfillment_order integer NOT NULL,
  order_id uuid NOT NULL,
  status varchar(32) NOT NULL,
  tracking_number varchar(255) NOT NULL,
  created_at bigint NOT NULL,
  shipping_refund_amount double precision,
  total_refund_amount double precision,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_fulfillments_status ON fulfillments USING btree (status);

CREATE INDEX idx_fulfillments_tracking_number ON fulfillments USING btree (tracking_number);