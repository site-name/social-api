CREATE TABLE IF NOT EXISTS voucher_customers (
  id varchar(36) NOT NULL PRIMARY KEY,
  voucher_id varchar(36) NOT NULL,
  customer_email varchar(128) NOT NULL
);

ALTER TABLE ONLY voucher_customers
    ADD CONSTRAINT voucher_customers_voucher_id_customer_email_key UNIQUE (voucher_id, customer_email);
