CREATE TABLE IF NOT EXISTS open_exchange_rates (
  id varchar(36) NOT NULL PRIMARY KEY,
  to_currency Currency NOT NULL,
  rate decimal(3,2),
  created_at bigint NOT NULL
);

ALTER TABLE ONLY open_exchange_rates
    ADD CONSTRAINT open_exchange_rates_to_currency_key UNIQUE (to_currency);

CREATE INDEX idx_open_exchange_to_currency ON open_exchange_rates USING btree (to_currency);