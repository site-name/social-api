CREATE TABLE IF NOT EXISTS voucher_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code varchar(10) NOT NULL,
  name varchar(255) NOT NULL,
  voucher_id uuid NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY voucher_translations
    ADD CONSTRAINT voucher_translations_language_code_voucher_id_key UNIQUE (language_code, voucher_id);
