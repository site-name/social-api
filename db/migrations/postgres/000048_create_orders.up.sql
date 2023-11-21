CREATE TABLE IF NOT EXISTS orders (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  status OrderStatus NOT NULL,
  user_id uuid,
  language_code LanguageCode NOT NULL,
  tracking_client_id uuid NOT NULL,
  billing_address_id uuid,
  shipping_address_id uuid,
  user_email varchar(128) NOT NULL,
  original_id uuid,
  origin OrderOrigin NOT NULL,
  currency varchar(200) NOT NULL,
  shipping_method_id uuid,
  collection_point_id uuid,
  shipping_method_name varchar(255),
  collection_point_name varchar(255),
  channel_id uuid NOT NULL,
  shipping_price_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  shipping_price_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  shipping_tax_rate decimal(5,4) NOT NULL DEFAULT 0.00,
  token uuid NOT NULL,
  checkout_token uuid NOT NULL,
  total_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_total_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_total_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_paid_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  voucher_id uuid,
  display_gross_prices boolean,
  customer_note text NOT NULL,
  weight_amount real NOT NULL,
  weight_unit varchar(10) NOT NULL,
  redirect_url text,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY orders
    ADD CONSTRAINT orders_token_key UNIQUE (token);

CREATE INDEX idx_orders_metadata ON orders USING btree (metadata);

CREATE INDEX idx_orders_private_metadata ON orders USING btree (private_metadata);

CREATE INDEX idx_orders_user_email ON orders USING btree (user_email);

CREATE INDEX idx_orders_user_email_lower_textpattern ON orders USING btree (lower((user_email)::text) text_pattern_ops);