CREATE TABLE IF NOT EXISTS voucher_customers (
  id varchar(36) NOT NULL PRIMARY KEY,
  voucher_code_id varchar(36) NOT NULL,
  customer_email varchar(128) NOT NULL
);

ALTER TABLE voucher_customers ADD CONSTRAINT fk_voucher_code_id FOREIGN KEY (voucher_code_id) REFERENCES voucher_codes(id) ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS idx_unique_together_voucher_customer ON voucher_customers USING btree (voucher_code_id, customer_email);
