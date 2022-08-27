CREATE TABLE IF NOT EXISTS shippingmethodtranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  shippingmethodid character varying(36),
  languagecode character varying(5),
  name character varying(100),
  description text
);

ALTER TABLE ONLY shippingmethodtranslations
    ADD CONSTRAINT shippingmethodtranslations_languagecode_shippingmethodid_key UNIQUE (languagecode, shippingmethodid);

CREATE INDEX idx_shipping_method_translations_name ON shippingmethodtranslations USING btree (name);

CREATE INDEX idx_shipping_method_translations_name_lower_textpattern ON shippingmethodtranslations USING btree (lower((name)::text) text_pattern_ops);
