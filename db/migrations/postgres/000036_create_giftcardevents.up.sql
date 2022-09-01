CREATE TABLE IF NOT EXISTS giftcardevents (
  id character varying(36) NOT NULL PRIMARY KEY,
  date bigint,
  type character varying(255),
  parameters jsonb,
  userid character varying(36),
  giftcardid character varying(36)
);

CREATE INDEX idx_giftcardevents_date ON giftcardevents USING btree (date);

