CREATE TABLE IF NOT EXISTS orders (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  status character varying(32),
  user_id uuid,
  language_code character varying(5),
  tracking_client_id uuid,
  billing_address_id uuid,
  shipping_address_id uuid,
  user_email character varying(128),
  original_id uuid,
  origin character varying(32),
  currency character varying(200),
  shipping_method_id uuid,
  collection_point_id uuid,
  shipping_method_name character varying(255),
  collection_point_name character varying(255),
  channel_id uuid,
  shipping_price_net_amount double precision,
  shipping_price_gross_amount double precision,
  shipping_tax_rate double precision,
  token uuid,
  checkout_token uuid,
  total_net_amount double precision,
  undiscounted_total_net_amount double precision,
  total_gross_amount double precision,
  undiscounted_total_gross_amount double precision,
  total_paid_amount double precision,
  voucher_id uuid,
  display_gross_prices boolean,
  customer_note text,
  weight_amount real,
  weight_unit text,
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