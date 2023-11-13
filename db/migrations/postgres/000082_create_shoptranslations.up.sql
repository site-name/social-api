CREATE TABLE IF NOT EXISTS shop_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(5),
  name character varying(110),
  description character varying(110),
  created_at bigint,
  updated_at bigint
);
