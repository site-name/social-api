CREATE TABLE IF NOT EXISTS voucher_customers (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  voucher_id uuid NOT NULL,
  customer_email character varying(128) NOT NULL
);

ALTER TABLE ONLY voucher_customers
    ADD CONSTRAINT voucher_customers_voucher_id_customer_email_key UNIQUE (voucher_id, customer_email);
