CREATE TABLE IF NOT EXISTS shop_translations (
  id varchar(36) NOT NULL PRIMARY KEY,
  language_code language_code NOT NULL,
  name varchar(110) NOT NULL,
  description varchar(110) NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL
);
