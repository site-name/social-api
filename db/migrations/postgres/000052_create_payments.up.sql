CREATE TABLE IF NOT EXISTS payments (
  id varchar(36) NOT NULL PRIMARY KEY,
  gateway varchar(255) NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  to_confirm boolean NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  charge_status payment_charge_status NOT NULL,
  token varchar(512) NOT NULL,
  total decimal(12,3) NOT NULL DEFAULT 0.00,
  captured_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  currency Currency NOT NULL,
  checkout_id varchar(36),
  order_id varchar(36),
  billing_email varchar(128) NOT NULL,
  billing_first_name varchar(256) NOT NULL,
  billing_last_name varchar(256) NOT NULL,
  billing_company_name varchar(256) NOT NULL,
  billing_address1 varchar(256) NOT NULL,
  billing_address2 varchar(256) NOT NULL,
  billing_city varchar(256) NOT NULL,
  billing_city_area varchar(128) NOT NULL,
  billing_postal_code varchar(20) NOT NULL,
  billing_country_code country_code NOT NULL,
  billing_country_area varchar(256) NOT NULL,
  cc_first_digits varchar(6) NOT NULL,
  cc_last_digits varchar(4) NOT NULL,
  cc_brand varchar(40) NOT NULL,
  cc_exp_month integer,
  cc_exp_year integer,
  payment_method_type varchar(256) NOT NULL,
  customer_ip_address varchar(39),
  extra_data text NOT NULL,
  return_url varchar(200),
  psp_reference varchar(512),
  store_payment_method store_payment_method NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_payments_charge_status ON payments USING btree (charge_status);

CREATE INDEX idx_payments_is_active ON payments USING btree (is_active);

CREATE INDEX idx_payments_metadata ON payments USING btree (metadata);

CREATE INDEX idx_payments_order_id ON payments USING btree (order_id);

CREATE INDEX idx_payments_private_metadata ON payments USING btree (private_metadata);

CREATE INDEX idx_payments_psp_reference ON payments USING btree (psp_reference);