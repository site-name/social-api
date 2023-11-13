CREATE TABLE IF NOT EXISTS payments (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  gateway character varying(255),
  is_active boolean,
  to_confirm boolean,
  created_at bigint,
  updated_at bigint,
  charge_status character varying(20),
  token character varying(512),
  total double precision,
  captured_amount double precision,
  currency character varying(3),
  checkout_id uuid,
  order_id uuid,
  billing_email character varying(128),
  billing_first_name character varying(256),
  billing_last_name character varying(256),
  billing_company_name character varying(256),
  billing_address1 character varying(256),
  billing_address2 character varying(256),
  billing_city character varying(256),
  billing_city_area character varying(128),
  billing_postal_code character varying(20),
  billing_country_code character varying(5),
  billing_country_area character varying(256),
  cc_first_digits character varying(6),
  cc_last_digits character varying(4),
  cc_brand character varying(40),
  cc_exp_month integer,
  cc_exp_year integer,
  payment_method_type character varying(256),
  customer_ip_address character varying(39),
  extra_data text,
  return_url character varying(200),
  psp_reference character varying(512),
  store_payment_method character varying(11),
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX idx_payments_charge_status ON payments USING btree (charge_status);

CREATE INDEX idx_payments_is_active ON payments USING btree (is_active);

CREATE INDEX idx_payments_metadata ON payments USING btree (metadata);

CREATE INDEX idx_payments_order_id ON payments USING btree (order_id);

CREATE INDEX idx_payments_private_metadata ON payments USING btree (private_metadata);

CREATE INDEX idx_payments_psp_reference ON payments USING btree (psp_reference);