CREATE TABLE IF NOT EXISTS plugin_configurations (
  id varchar(36) NOT NULL PRIMARY KEY,
  identifier varchar(128) NOT NULL,
  name varchar(128) NOT NULL,
  channel_id varchar(36) nOT NULL,
  description varchar(1000) NOT NULL,
  active boolean NOT NULL,
  configuration jsonb
);

ALTER TABLE ONLY plugin_configurations
    ADD CONSTRAINT plugin_configurations_identifier_channel_id_key UNIQUE (identifier, channel_id);

CREATE INDEX idx_plugin_configurations_identifier ON plugin_configurations USING btree (identifier);

CREATE INDEX idx_plugin_configurations_lower_textpattern_name ON plugin_configurations USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_plugin_configurations_name ON plugin_configurations USING btree (name);