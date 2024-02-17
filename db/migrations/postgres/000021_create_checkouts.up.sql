CREATE TABLE IF NOT EXISTS checkouts (
  token varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  user_id varchar(36),
  email text NOT NULL,
  quantity integer NOT NULL,
  channel_id varchar(36) NOT NULL,
  billing_address_id varchar(36),
  shipping_address_id varchar(36),
  shipping_method_id varchar(36),
  collection_point_id varchar(36),
  note text NOT NULL,
  currency Currency NOT NULL,
  country country_code NOT NULL,
  discount_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  discount_name varchar(255),
  translated_discount_name varchar(255),
  voucher_code varchar(12),
  redirect_url text,
  tracking_code varchar(255),
  language_code language_code NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_checkouts_billing_address_id ON checkouts USING btree (billing_address_id);

CREATE INDEX idx_checkouts_channel_id ON checkouts USING btree (channel_id);

CREATE INDEX idx_checkouts_shipping_address_id ON checkouts USING btree (shipping_address_id);

CREATE INDEX idx_checkouts_shipping_method_id ON checkouts USING btree (shipping_method_id);

CREATE INDEX idx_checkouts_token ON checkouts USING btree (token);

CREATE INDEX idx_checkouts_user_id ON checkouts USING btree (user_id);