CREATE TABLE IF NOT EXISTS pluginconfigurations (
  id character varying(36) NOT NULL PRIMARY KEY,
  identifier character varying(128),
  name character varying(128),
  channelid character varying(36),
  description character varying(1000),
  active boolean,
  configuration text
);

ALTER TABLE ONLY pluginconfigurations
    ADD CONSTRAINT pluginconfigurations_identifier_channelid_key UNIQUE (identifier, channelid);

CREATE INDEX idx_plugin_configurations_identifier ON pluginconfigurations USING btree (identifier);

CREATE INDEX idx_plugin_configurations_lower_textpattern_name ON pluginconfigurations USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_plugin_configurations_name ON pluginconfigurations USING btree (name);
