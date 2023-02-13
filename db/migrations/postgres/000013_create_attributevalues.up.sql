CREATE TABLE IF NOT EXISTS attributevalues (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  value character varying(9),
  slug character varying(255),
  fileurl character varying(200),
  contenttype character varying(50),
  attributeid character varying(36),
  richtext text,
  "boolean" boolean,
  datetime timestamp with time zone,
  sortorder integer
);

ALTER TABLE ONLY attributevalues
    ADD CONSTRAINT attributevalues_slug_attributeid_key UNIQUE (slug, attributeid);

CREATE INDEX idx_attributevalues_name ON attributevalues USING btree (name);

CREATE INDEX idx_attributevalues_name_lower_textpattern ON attributevalues USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_attributevalues_slug ON attributevalues USING btree (slug);

ALTER TABLE ONLY attributevalues
    ADD CONSTRAINT fk_attributevalues_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;
