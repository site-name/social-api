CREATE TABLE IF NOT EXISTS collections (
  id character varying(36) NOT NULL PRIMARY KEY,
  shopid character varying(36),
  name character varying(250),
  slug character varying(255),
  backgroundimage character varying(200),
  backgroundimagealt character varying(128),
  description text,
  metadata jsonb,
  privatemetadata jsonb,
  seotitle character varying(70),
  seodescription character varying(300)
);

CREATE INDEX idx_collections_name ON collections USING btree (name);

CREATE INDEX idx_collections_name_lower_textpattern ON collections USING btree (lower((name)::text) text_pattern_ops);

ALTER TABLE ONLY collections
    ADD CONSTRAINT fk_collections_shops FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;
