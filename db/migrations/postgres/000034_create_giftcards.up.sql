CREATE TABLE IF NOT EXISTS giftcards (
  id varchar(36) NOT NULL PRIMARY KEY,
  code varchar(40) NOT NULL,
  created_by_id varchar(36),
  used_by_id varchar(36),
  created_by_email varchar(128),
  used_by_email varchar(128),
  created_at bigint NOT NULL,
  expiry_date timestamp with time zone,
  tag varchar(255),
  product_id varchar(36),
  last_used_on bigint,
  is_active boolean,
  currency Currency NOT NULL,
  initial_balance_amount decimal(12,3),
  current_balance_amount decimal(12,3),
  app_id varchar(36),
  fulfillment_line_id varchar(36),
  search_vector tsvector,
  search_index_dirty boolean,
  annotations jsonb, -- This field is used to store additional information about the gift card, such as the sender's name, the recipient's name, etc.
  metadata jsonb,
  private_metadata jsonb
);

CREATE TABLE IF NOT EXISTS giftcard_tags (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS giftcard_tag_giftcards (
  id varchar(36) NOT NULL PRIMARY KEY,
  giftcard_id varchar(36) NOT NULL,
  tag_id varchar(36) NOT NULL
);

ALTER TABLE ONLY giftcard_tag_giftcards
    ADD CONSTRAINT fk_giftcard_tag_giftcards_giftcards FOREIGN KEY (giftcard_id) REFERENCES giftcards(id) ON DELETE CASCADE;
ALTER TABLE ONLY giftcard_tag_giftcards
    ADD CONSTRAINT fk_giftcard_tag_giftcards_giftcard_tags FOREIGN KEY (tag_id) REFERENCES giftcard_tags(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX unique_giftcard_id_tag_id ON giftcard_tag_giftcards (giftcard_id, tag_id);

ALTER TABLE ONLY giftcards
    ADD CONSTRAINT giftcards_code_key UNIQUE (code);

CREATE INDEX idx_giftcards_code ON giftcards USING btree (code);
CREATE INDEX idx_giftcards_metadata ON giftcards USING btree (metadata);
CREATE INDEX idx_giftcards_private_metadata ON giftcards USING btree (private_metadata);
CREATE INDEX idx_giftcards_tag ON giftcards USING btree (tag);

CREATE UNIQUE INDEX unique_giftcard_tags_name ON giftcard_tags (name);
