CREATE TABLE IF NOT EXISTS fulfillments (
  id varchar(36) NOT NULL PRIMARY KEY,
  fulfillment_order integer NOT NULL,
  order_id varchar(36) NOT NULL,
  status fulfillment_status NOT NULL,
  tracking_number varchar(255) NOT NULL,
  created_at bigint NOT NULL,
  shipping_refund_amount decimal(12,3),
  total_refund_amount decimal(12,3),
  metadata jsonb,
  private_metadata jsonb
);


CREATE INDEX idx_fulfillments_status ON fulfillments USING btree (status);

CREATE INDEX idx_fulfillments_tracking_number ON fulfillments USING btree (tracking_number);