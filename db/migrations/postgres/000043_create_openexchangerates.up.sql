CREATE TABLE IF NOT EXISTS open_exchange_rates (
  id character varying(36) NOT NULL PRIMARY KEY,
  tocurrency character varying(3),
  rate double precision
);

ALTER TABLE ONLY open_exchange_rates
    ADD CONSTRAINT open_exchange_rates_tocurrency_key UNIQUE (tocurrency);

CREATE INDEX idx_openexchange_to_currency ON open_exchange_rates USING btree (tocurrency);
