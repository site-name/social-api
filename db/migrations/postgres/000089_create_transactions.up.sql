CREATE TABLE IF NOT EXISTS transactions (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  paymentid character varying(36),
  token character varying(512),
  kind character varying(25),
  issuccess boolean,
  actionrequired boolean,
  actionrequireddata text,
  currency character varying(3),
  amount double precision,
  error character varying(256),
  customerid character varying(256),
  gatewayresponse text,
  alreadyprocessed boolean
);
ALTER TABLE ONLY transactions
    ADD CONSTRAINT fk_transactions_payments FOREIGN KEY (paymentid) REFERENCES payments(id);
