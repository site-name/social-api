CREATE TABLE IF NOT EXISTS invoices (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid,
  number character varying(255) NOT NULL,
  created_at bigint NOT NULL,
  external_url character varying(2048) NOT NULL,
  status character varying(50) NOT NULL,
  message character varying(255) NOT NULL,
  updated_at bigint NOT NULL,
  invoice_file uuid, -- the file id
  metadata jsonb,
  private_metadata jsonb
);