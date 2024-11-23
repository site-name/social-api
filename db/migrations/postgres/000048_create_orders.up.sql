CREATE TABLE IF NOT EXISTS orders (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  expired_at bigint,
  status order_status NOT NULL,
  authorize_status order_authorize_status NOT NULL DEFAULT 'none',
  charge_status order_charge_status NOT NULL DEFAULT 'none',
  user_id varchar(36),
  language_code language_code NOT NULL,
  tracking_client_id varchar(36) NOT NULL,
  billing_address_id varchar(36),
  shipping_address_id varchar(36),
  user_email varchar(128) NOT NULL,
  original_id varchar(36),
  origin order_origin,
  currency Currency NOT NULL,
  shipping_method_id varchar(36),
  collection_point_id varchar(36),
  shipping_method_name varchar(255),
  collection_point_name varchar(255),
  channel_id varchar(36) NOT NULL,
  shipping_price_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  shipping_price_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  base_shipping_price_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_base_shipping_price_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  shipping_tax_rate decimal(5,4) NOT NULL DEFAULT 0.00,
  shipping_tax_class_id varchar(36),
  shipping_tax_class_name varchar(255),
  shipping_tax_class_private_metadata jsonb,
  shipping_tax_class_metadata jsonb,
  token varchar(36) NOT NULL,
  checkout_token varchar(36) NOT NULL,
  total_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_total_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_total_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_paid_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_charged_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_authorized_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  subtotal_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  subtotal_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  voucher_code varchar(255),
  voucher_id varchar(36),
  display_gross_prices boolean,
  customer_note text NOT NULL,
  weight_amount real NOT NULL,
  weight_unit varchar(10) NOT NULL,
  redirect_url varchar(1000),
  search_document text,
  search_vector tsvector,
  tax_error varchar(255),
  tax_exemption boolean NOT NULL DEFAULT false,
  should_refresh_prices boolean NOT NULL DEFAULT true,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY orders
    ADD CONSTRAINT orders_token_key UNIQUE (token);

CREATE INDEX idx_orders_metadata ON orders USING btree (metadata);

CREATE INDEX idx_orders_private_metadata ON orders USING btree (private_metadata);

CREATE INDEX idx_orders_user_email ON orders USING btree (user_email);

CREATE INDEX idx_orders_user_email_lower_textpattern ON orders USING btree (lower((user_email)::text) text_pattern_ops);