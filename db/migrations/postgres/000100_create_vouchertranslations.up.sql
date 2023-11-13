CREATE TABLE IF NOT EXISTS voucher_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(10),
  name character varying(255),
  voucher_id uuid,
  created_at bigint
);

ALTER TABLE ONLY voucher_translations
    ADD CONSTRAINT voucher_translations_language_code_voucher_id_key UNIQUE (language_code, voucher_id);
