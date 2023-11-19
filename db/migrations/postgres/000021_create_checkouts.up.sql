CREATE TABLE IF NOT EXISTS checkouts (
  token uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  user_id uuid,
  email text NOT NULL,
  quantity integer NOT NULL,
  channel_id uuid NOT NULL,
  billing_address_id uuid,
  shipping_address_id uuid,
  shipping_method_id uuid,
  collection_point_id uuid,
  note text NOT NULL,
  currency text NOT NULL,
  country character varying(5) NOT NULL,
  discount_amount double precision,
  discount_name character varying(255),
  translated_discount_name character varying(255),
  voucher_code character varying(12),
  redirect_url text,
  tracking_code character varying(255),
  language_code text NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_checkouts_billing_address_id ON checkouts USING btree (billing_address_id);

CREATE INDEX idx_checkouts_channel_id ON checkouts USING btree (channel_id);

CREATE INDEX idx_checkouts_shipping_address_id ON checkouts USING btree (shipping_address_id);

CREATE INDEX idx_checkouts_shipping_method_id ON checkouts USING btree (shipping_method_id);

CREATE INDEX idx_checkouts_token ON checkouts USING btree (token);

CREATE INDEX idx_checkouts_user_id ON checkouts USING btree (user_id);