CREATE TABLE IF NOT EXISTS orders (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  status character varying(32),
  userid character varying(36),
  shopid character varying(36),
  languagecode character varying(5),
  trackingclientid character varying(36),
  billingaddressid character varying(36),
  shippingaddressid character varying(36),
  useremail character varying(128),
  originalid character varying(36),
  origin character varying(32),
  currency character varying(200),
  shippingmethodid character varying(36),
  collectionpointid character varying(36),
  shippingmethodname character varying(255),
  collectionpointname character varying(255),
  channelid character varying(36),
  shippingpricenetamount double precision,
  shippingpricegrossamount double precision,
  shippingtaxrate double precision,
  token character varying(36),
  checkouttoken character varying(36),
  totalnetamount double precision,
  undiscountedtotalnetamount double precision,
  totalgrossamount double precision,
  undiscountedtotalgrossamount double precision,
  totalpaidamount double precision,
  voucherid character varying(36),
  displaygrossprices boolean,
  customernote text,
  weightamount real,
  weightunit text,
  redirecturl text,
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY orders
    ADD CONSTRAINT orders_token_key UNIQUE (token);

CREATE INDEX idx_orders_metadata ON orders USING btree (metadata);

CREATE INDEX idx_orders_private_metadata ON orders USING btree (privatemetadata);

CREATE INDEX idx_orders_user_email ON orders USING btree (useremail);

CREATE INDEX idx_orders_user_email_lower_textpattern ON orders USING btree (lower((useremail)::text) text_pattern_ops);
