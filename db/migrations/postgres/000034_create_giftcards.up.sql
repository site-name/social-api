CREATE TABLE IF NOT EXISTS giftcards (
  id character varying(36) NOT NULL PRIMARY KEY,
  code character varying(40),
  createdbyid character varying(36),
  usedbyid character varying(36),
  createdbyemail character varying(128),
  usedbyemail character varying(128),
  createat bigint,
  startdate timestamp with time zone,
  expirydate timestamp with time zone,
  tag character varying(255),
  productid character varying(36),
  lastusedon bigint,
  isactive boolean,
  currency character varying(3),
  initialbalanceamount double precision,
  currentbalanceamount double precision,
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY giftcards
    ADD CONSTRAINT giftcards_code_key UNIQUE (code);

CREATE INDEX idx_giftcards_code ON giftcards USING btree (code);

CREATE INDEX idx_giftcards_metadata ON giftcards USING btree (metadata);

CREATE INDEX idx_giftcards_private_metadata ON giftcards USING btree (privatemetadata);

CREATE INDEX idx_giftcards_tag ON giftcards USING btree (tag);

ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_products FOREIGN KEY (productid) REFERENCES products(id);

ALTER TABLE ONLY giftcards
    ADD CONSTRAINT fk_giftcards_users FOREIGN KEY (createdbyid) REFERENCES users(id);
