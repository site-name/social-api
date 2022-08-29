CREATE TABLE IF NOT EXISTS payments (
  id character varying(36) NOT NULL PRIMARY KEY,
  gateway character varying(255),
  isactive boolean,
  toconfirm boolean,
  createat bigint,
  updateat bigint,
  chargestatus character varying(20),
  token character varying(512),
  total double precision,
  capturedamount double precision,
  currency character varying(3),
  checkoutid character varying(36),
  orderid character varying(36),
  billingemail character varying(128),
  billingfirstname character varying(256),
  billinglastname character varying(256),
  billingcompanyname character varying(256),
  billingaddress1 character varying(256),
  billingaddress2 character varying(256),
  billingcity character varying(256),
  billingcityarea character varying(128),
  billingpostalcode character varying(20),
  billingcountrycode character varying(5),
  billingcountryarea character varying(256),
  ccfirstdigits character varying(6),
  cclastdigits character varying(4),
  ccbrand character varying(40),
  ccexpmonth integer,
  ccexpyear integer,
  paymentmethodtype character varying(256),
  customeripaddress character varying(39),
  extradata text,
  returnurl character varying(200),
  pspreference character varying(512),
  storepaymentmethod character varying(11),
  metadata jsonb,
  privatemetadata jsonb
);

CREATE INDEX idx_payments_charge_status ON payments USING btree (chargestatus);

CREATE INDEX idx_payments_is_active ON payments USING btree (isactive);

CREATE INDEX idx_payments_metadata ON payments USING btree (metadata);

CREATE INDEX idx_payments_order_id ON payments USING btree (orderid);

CREATE INDEX idx_payments_private_metadata ON payments USING btree (privatemetadata);

CREATE INDEX idx_payments_psp_reference ON payments USING btree (pspreference);

ALTER TABLE ONLY payments
    ADD CONSTRAINT fk_payments_checkouts FOREIGN KEY (checkoutid) REFERENCES checkouts(token);

ALTER TABLE ONLY payments
    ADD CONSTRAINT fk_payments_orders FOREIGN KEY (orderid) REFERENCES orders(id);
