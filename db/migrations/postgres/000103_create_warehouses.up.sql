CREATE TABLE IF NOT EXISTS warehouses (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(250) NOT NULL,
  slug varchar(255) NOT NULL,
  address_id varchar(36),
  email varchar(128) NOT NULL,
  click_and_collect_option warehouse_click_and_collect_option NOT NULL,
  is_private boolean,
  created_at bigint NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY warehouses
    ADD CONSTRAINT warehouses_slug_key UNIQUE (slug);

CREATE INDEX idx_warehouses_email ON warehouses USING btree (email);

CREATE INDEX idx_warehouses_email_lower_text_pattern ON warehouses USING btree (lower((email)::text) text_pattern_ops);
