CREATE TABLE IF NOT EXISTS attributes (
  id varchar(36) NOT NULL PRIMARY KEY,
  slug varchar(255) NOT NULL,
  name varchar(250) NOT NULL,
  type attribute_type NOT NULL, -- enum
  input_type attribute_input_type NOT NULL, -- enum
  entity_type attribute_entity_type, -- enum
  unit varchar(100),
  value_required boolean NOT NULL,
  is_variant_only boolean NOT NULL,
  visible_in_storefront boolean NOT NULL,
  filterable_in_storefront boolean NOT NULL,
  filterable_in_dashboard boolean NOT NULL,
  storefront_search_position integer NOT NULL,
  available_in_grid boolean NOT NULL,
  metadata jsonb,
  private_metadata jsonb,
  max_sfort_order integer
);


ALTER TABLE ONLY attributes
    ADD CONSTRAINT attributes_slug_key UNIQUE (slug);

CREATE INDEX idx_attributes_name ON attributes USING btree (name);

CREATE INDEX idx_attributes_name_lower_textpattern ON attributes USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_attributes_slug ON attributes USING btree (slug);