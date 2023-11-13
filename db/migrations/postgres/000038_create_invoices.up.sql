CREATE TABLE IF NOT EXISTS invoices (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid,
  number character varying(255),
  created_at bigint,
  external_url character varying(2048),
  status character varying(50),
  message character varying(255),
  updated_at bigint,
  invoice_file text,
  metadata jsonb,
  private_metadata jsonb
);