CREATE TABLE IF NOT EXISTS open_exchange_rates (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  to_currency varchar(3) NOT NULL,
  rate double precision
);

ALTER TABLE ONLY open_exchange_rates
    ADD CONSTRAINT open_exchange_rates_to_currency_key UNIQUE (to_currency);

CREATE INDEX idx_open_exchange_to_currency ON open_exchange_rates USING btree (to_currency);