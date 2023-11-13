CREATE TABLE IF NOT EXISTS plugin_configurations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  identifier character varying(128),
  name character varying(128),
  channel_id uuid,
  description character varying(1000),
  active boolean,
  configuration text
);

ALTER TABLE ONLY plugin_configurations
    ADD CONSTRAINT plugin_configurations_identifier_channel_id_key UNIQUE (identifier, channel_id);

CREATE INDEX idx_plugin_configurations_identifier ON plugin_configurations USING btree (identifier);

CREATE INDEX idx_plugin_configurations_lower_textpattern_name ON plugin_configurations USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_plugin_configurations_name ON plugin_configurations USING btree (name);