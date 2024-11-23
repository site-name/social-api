CREATE TABLE IF NOT EXISTS voucher_codes (
  id varchar(36) NOT NULL PRIMARY KEY,
  voucher_id varchar(36) NOT NULL,
  code varchar(255) NOT NULL,
  used integer NOT NULL DEFAULT 0,
  is_active boolean NOT NULL DEFAULT false,
  created_at bigint NOT NULL
);

ALTER TABLE voucher_codes ADD CONSTRAINT fk_voucher_id FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;
CREATE INDEX idx_voucher_code ON voucher_codes USING btree (code);
