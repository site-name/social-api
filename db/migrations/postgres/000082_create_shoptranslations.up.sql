CREATE TABLE IF NOT EXISTS shop_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(5) NOT NULL,
  name character varying(110) NOT NULL,
  description character varying(110) NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL
);
