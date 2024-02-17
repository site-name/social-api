CREATE TABLE IF NOT EXISTS voucher_translations (
  id varchar(36) NOT NULL PRIMARY KEY,
  language_code language_code NOT NULL,
  name varchar(255) NOT NULL,
  voucher_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY voucher_translations
    ADD CONSTRAINT voucher_translations_language_code_voucher_id_key UNIQUE (language_code, voucher_id);
