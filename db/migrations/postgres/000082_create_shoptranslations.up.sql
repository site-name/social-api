CREATE TABLE IF NOT EXISTS shop_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code LanguageCode NOT NULL,
  name varchar(110) NOT NULL,
  description varchar(110) NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL
);
