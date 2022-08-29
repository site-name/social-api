CREATE TABLE IF NOT EXISTS attributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  slug character varying(255),
  name character varying(250),
  type character varying(50),
  inputtype character varying(50),
  entitytype character varying(50),
  unit character varying(100),
  valuerequired boolean,
  isvariantonly boolean,
  visibleinstorefront boolean,
  filterableinstorefront boolean,
  filterableindashboard boolean,
  storefrontsearchposition integer,
  availableingrid boolean,
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY attributes
    ADD CONSTRAINT attributes_slug_key UNIQUE (slug);

CREATE INDEX idx_attributes_name ON attributes USING btree (name);

CREATE INDEX idx_attributes_name_lower_textpattern ON attributes USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_attributes_slug ON attributes USING btree (slug);
