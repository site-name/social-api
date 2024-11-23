CREATE TABLE IF NOT EXISTS voucher_customers (
  id varchar(36) NOT NULL PRIMARY KEY,
  voucher_code_id varchar(36) NOT NULL,
  customer_email varchar(128) NOT NULL
);

ALTER TABLE ONLY voucher_customers
    ADD CONSTRAINT voucher_customers_voucher_code_id_customer_email_key UNIQUE (voucher_code_id, customer_email);
CREATE INDEX IF NOT EXISTS idx_voucher_code_id ON voucher_customers USING btree (voucher_code_id);
