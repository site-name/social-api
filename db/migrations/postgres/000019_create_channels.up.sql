CREATE TABLE IF NOT EXISTS channels (
  id character varying(36) NOT NULL PRIMARY KEY,
  shopid character varying(36),
  name character varying(250),
  isactive boolean,
  slug character varying(255),
  currency text,
  defaultcountry character varying(5)
);

ALTER TABLE ONLY channels
    ADD CONSTRAINT channels_slug_key UNIQUE (slug);

CREATE INDEX IF NOT EXISTS idx_channels_currency ON channels USING btree (currency);

CREATE INDEX IF NOT EXISTS idx_channels_name ON channels USING btree (name);

CREATE INDEX IF NOT EXISTS idx_channels_isactive ON channels USING btree (isactive);

CREATE INDEX IF NOT EXISTS idx_channels_name_lower_textpattern ON channels USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_channels_slug ON channels USING btree (slug);
