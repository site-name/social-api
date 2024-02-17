CREATE TABLE IF NOT EXISTS giftcards (
  id varchar(36) NOT NULL PRIMARY KEY,
  code varchar(40) NOT NULL,
  created_by_id varchar(36),
  used_by_id varchar(36),
  created_by_email varchar(128),
  used_by_email varchar(128),
  created_at bigint NOT NULL,
  start_date timestamp with time zone,
  expiry_date timestamp with time zone,
  tag varchar(255),
  product_id varchar(36),
  last_used_on bigint,
  is_active boolean,
  currency Currency NOT NULL,
  initial_balance_amount decimal(12,3),
  current_balance_amount decimal(12,3),
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY giftcards
    ADD CONSTRAINT giftcards_code_key UNIQUE (code);

CREATE INDEX idx_giftcards_code ON giftcards USING btree (code);

CREATE INDEX idx_giftcards_metadata ON giftcards USING btree (metadata);

CREATE INDEX idx_giftcards_private_metadata ON giftcards USING btree (private_metadata);

CREATE INDEX idx_giftcards_tag ON giftcards USING btree (tag);