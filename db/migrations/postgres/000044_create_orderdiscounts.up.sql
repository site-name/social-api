CREATE TABLE IF NOT EXISTS orderdiscounts (
  id character varying(36) NOT NULL PRIMARY KEY,
  orderid character varying(36),
  type character varying(10),
  valuetype character varying(10),
  value double precision,
  amountvalue double precision,
  currency text,
  name character varying(255),
  translatedname character varying(255),
  reason text
);

CREATE INDEX idx_order_discounts_name ON orderdiscounts USING btree (name);

CREATE INDEX idx_order_discounts_translated_name ON orderdiscounts USING btree (translatedname);

CREATE INDEX idx_order_discounts_name_lower_textpattern ON orderdiscounts USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_order_discounts_translated_name_lower_textpattern ON orderdiscounts USING btree (lower((translatedname)::text) text_pattern_ops);
