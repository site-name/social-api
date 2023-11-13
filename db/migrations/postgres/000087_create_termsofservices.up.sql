CREATE TABLE IF NOT EXISTS terms_of_services (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  user_id uuid,
  text character varying(65535)
);