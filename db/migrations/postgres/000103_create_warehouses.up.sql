CREATE TABLE IF NOT EXISTS warehouses (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  slug character varying(255),
  addressid character varying(36),
  email character varying(128),
  clickandcollectoption character varying(30),
  isprivate boolean,
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY warehouses
    ADD CONSTRAINT warehouses_slug_key UNIQUE (slug);

CREATE INDEX idx_warehouses_email ON warehouses USING btree (email);

CREATE INDEX idx_warehouses_email_lower_textpattern ON warehouses USING btree (lower((email)::text) text_pattern_ops);
