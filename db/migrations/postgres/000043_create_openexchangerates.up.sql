CREATE TABLE IF NOT EXISTS openexchangerates (
  id character varying(36) NOT NULL PRIMARY KEY,
  tocurrency character varying(3),
  rate double precision
);

ALTER TABLE ONLY openexchangerates
    ADD CONSTRAINT openexchangerates_tocurrency_key UNIQUE (tocurrency);

CREATE INDEX idx_openexchange_to_currency ON openexchangerates USING btree (tocurrency);
