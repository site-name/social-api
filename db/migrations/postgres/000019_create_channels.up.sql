CREATE TABLE IF NOT EXISTS channels (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(250) NOT NULL,
  is_active boolean NOT NULL,
  slug varchar(255) NOT NULL,
  currency text NOT NULL,
  default_country country_code NOT NULL
);

ALTER TABLE ONLY channels
    ADD CONSTRAINT channels_slug_key UNIQUE (slug);

CREATE INDEX IF NOT EXISTS idx_channels_currency ON channels USING btree (currency);

CREATE INDEX IF NOT EXISTS idx_channels_name ON channels USING btree (name);

CREATE INDEX IF NOT EXISTS idx_channels_is_active ON channels USING btree (is_active);

CREATE INDEX IF NOT EXISTS idx_channels_name_lower_textpattern ON channels USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_channels_slug ON channels USING btree (slug);