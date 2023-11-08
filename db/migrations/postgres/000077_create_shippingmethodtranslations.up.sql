CREATE TABLE IF NOT EXISTS shipping_method_translations (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingmethodid character varying(36),
  languagecode character varying(5),
  name character varying(100),
  description text
);

ALTER TABLE ONLY shipping_method_translations
    ADD CONSTRAINT shipping_method_translations_languagecode_shippingmethodid_key UNIQUE (languagecode, shippingmethodid);

CREATE INDEX idx_shipping_method_translations_name ON shipping_method_translations USING btree (name);

CREATE INDEX idx_shipping_method_translations_name_lower_textpattern ON shipping_method_translations USING btree (lower((name)::text) text_pattern_ops);
