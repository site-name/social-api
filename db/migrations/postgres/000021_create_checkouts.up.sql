CREATE TABLE IF NOT EXISTS checkouts (
  token character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  updateat bigint,
  userid character varying(36),
  email text,
  quantity integer,
  channelid character varying(36),
  billingaddressid character varying(36),
  shippingaddressid character varying(36),
  shippingmethodid character varying(36),
  collectionpointid character varying(36),
  note text,
  currency text,
  country character varying(5),
  discountamount double precision,
  discountname character varying(255),
  translateddiscountname character varying(255),
  vouchercode character varying(12),
  redirecturl text,
  trackingcode character varying(255),
  languagecode text,
  metadata jsonb,
  privatemetadata jsonb
);

CREATE INDEX idx_checkouts_billing_address_id ON checkouts USING btree (billingaddressid);

CREATE INDEX idx_checkouts_channelid ON checkouts USING btree (channelid);

CREATE INDEX idx_checkouts_shipping_address_id ON checkouts USING btree (shippingaddressid);

CREATE INDEX idx_checkouts_shipping_method_id ON checkouts USING btree (shippingmethodid);

CREATE INDEX idx_checkouts_token ON checkouts USING btree (token);

CREATE INDEX idx_checkouts_userid ON checkouts USING btree (userid);

