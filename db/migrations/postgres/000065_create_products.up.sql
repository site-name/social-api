CREATE TABLE IF NOT EXISTS products (
  id character varying(36) NOT NULL PRIMARY KEY,
  producttypeid character varying(36),
  name character varying(250),
  slug character varying(255),
  description text,
  descriptionplaintext text,
  categoryid character varying(36),
  createat bigint,
  updateat bigint,
  chargetaxes boolean,
  weight real,
  weightunit text,
  defaultvariantid character varying(36),
  rating real,
  metadata text,
  privatemetadata text,
  seotitle character varying(70),
  seodescription character varying(300)
);

ALTER TABLE ONLY products
    ADD CONSTRAINT products_name_key UNIQUE (name);

ALTER TABLE ONLY products
    ADD CONSTRAINT products_slug_key UNIQUE (slug);

CREATE INDEX idx_products_metadata ON products USING btree (metadata);

CREATE INDEX idx_products_name ON products USING btree (name);

CREATE INDEX idx_products_name_lower_textpattern ON products USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_products_private_metadata ON products USING btree (privatemetadata);

CREATE INDEX idx_products_slug ON products USING btree (slug);
