CREATE TABLE IF NOT EXISTS invoices (
  id varchar(36) NOT NULL PRIMARY KEY,
  order_id varchar(36),
  number varchar(255) NOT NULL,
  created_at bigint NOT NULL,
  external_url varchar(2048) NOT NULL,
  status varchar(50) NOT NULL,
  message varchar(255) NOT NULL,
  updated_at bigint NOT NULL,
  invoice_file varchar(36), -- the file id
  metadata jsonb,
  private_metadata jsonb
);
