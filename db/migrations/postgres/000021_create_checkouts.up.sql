CREATE TABLE IF NOT EXISTS checkouts (
  token varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  completing_started_at bigint,
  last_transaction_modified_at bigint,
  automatically_refundable boolean NOT NULL default false,
  user_id varchar(36),
  email text NOT NULL,
  quantity integer NOT NULL,
  channel_id varchar(36) NOT NULL,
  billing_address_id varchar(36),
  shipping_address_id varchar(36),
  shipping_method_id varchar(36),
  collection_point_id varchar(36), -- foreig key to warehouse
  total_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  base_total_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  subtotal_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  subtotal_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  base_subtotal_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  shipping_price_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  shipping_price_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  shipping_tax_rate decimal(5,4) NOT NULL DEFAULT 0.00,
  authorize_status checkout_authorize_status NOT NULL DEFAULT 'none',
  charge_status checkout_charge_status NOT NULL DEFAULT 'none',
  price_expiration bigint, -- default to timestamp now
  note text NOT NULL,
  currency Currency NOT NULL,
  country country_code NOT NULL,
  discount_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  discount_name varchar(255),
  translated_discount_name varchar(255),
  voucher_code varchar(12),
  is_voucher_usage_increased boolean NOT NULL DEFAULT false,
  redirect_url varchar(1000),
  tracking_code varchar(255),
  language_code language_code NOT NULL,
  tax_exemption boolean NOT NULL DEFAULT false,
  tax_error varchar(255),
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_checkouts_billing_address_id ON checkouts USING btree (billing_address_id);
CREATE INDEX idx_checkouts_channel_id ON checkouts USING btree (channel_id);
CREATE INDEX idx_checkouts_shipping_address_id ON checkouts USING btree (shipping_address_id);
CREATE INDEX idx_checkouts_shipping_method_id ON checkouts USING btree (shipping_method_id);
CREATE INDEX idx_checkouts_token ON checkouts USING btree (token);
CREATE INDEX idx_checkouts_user_id ON checkouts USING btree (user_id);
CREATE INDEX idx_checkouts_updated_at ON checkouts USING btree (updated_at);
