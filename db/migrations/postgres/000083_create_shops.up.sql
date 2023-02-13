CREATE TABLE IF NOT EXISTS shops (
  id character varying(36) NOT NULL PRIMARY KEY,
  ownerid character varying(36),
  createat bigint,
  updateat bigint,
  name character varying(100),
  headertext character varying(100)
  description character varying(200),
  topmenuid character varying(36),
  bottommenuid character varying(36),
  includetaxesinprice boolean,
  displaygrossprices boolean,
  chargetaxesonshipping boolean,
  trackinventorybydefault boolean,
  defaultweightunit character varying(10),
  automaticfulfillmentdigitalproducts boolean,
  defaultdigitalmaxdownloads integer,
  defaultdigitalurlvaliddays integer,
  addressid character varying(36),
  companyaddressid character varying(36),
  defaultmailsendername character varying(78),
  defaultmailsenderaddress text,
  customersetpasswordurl text,
  automaticallyconfirmallneworders boolean,
  fulfillmentautoapprove boolean,
  fulfillmentallowunpaid boolean,
  giftcardexpirytype character varying(32),
  giftcardexpiryperiodtype character varying(32),
  giftcardexpiryperiod integer,
  automaticallyfulfillnonshippablegiftcard boolean
);

CREATE INDEX idx_shops_description ON shops USING btree (description);

CREATE INDEX idx_shops_description_lower_textpattern ON shops USING btree (lower((description)::text) text_pattern_ops);

CREATE INDEX idx_shops_name ON shops USING btree (name);

CREATE INDEX idx_shops_name_lower_textpattern ON shops USING btree (lower((name)::text) text_pattern_ops);

ALTER TABLE ONLY shops
  ADD CONSTRAINT fk_shops_bottommenuid FOREIGN KEY (bottommenuid) REFERENCES menus(id);

ALTER TABLE ONLY shops
    ADD CONSTRAINT fk_shops_addressid FOREIGN KEY (addressid) REFERENCES addresses(id);
ALTER TABLE ONLY shops
    ADD CONSTRAINT fk_shops_companyaddressid FOREIGN KEY (companyaddressid) REFERENCES addresses(id);

ALTER TABLE ONLY shops
    ADD CONSTRAINT fk_shops_topmenuid FOREIGN KEY (topmenuid) REFERENCES menus(id);
ALTER TABLE ONLY shops
    ADD CONSTRAINT fk_shops_users FOREIGN KEY (ownerid) REFERENCES users(id) ON DELETE CASCADE;
