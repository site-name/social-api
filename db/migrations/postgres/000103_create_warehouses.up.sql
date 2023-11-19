CREATE TABLE IF NOT EXISTS warehouses (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(250) NOT NULL,
  slug character varying(255) NOT NULL,
  address_id uuid,
  email character varying(128) NOT NULL,
  click_and_collect_option character varying(30) NOT NULL,
  is_private boolean,
  created_at bigint NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY warehouses
    ADD CONSTRAINT warehouses_slug_key UNIQUE (slug);

CREATE INDEX idx_warehouses_email ON warehouses USING btree (email);

CREATE INDEX idx_warehouses_email_lower_text_pattern ON warehouses USING btree (lower((email)::text) text_pattern_ops);
