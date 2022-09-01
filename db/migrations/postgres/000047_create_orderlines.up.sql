CREATE TABLE IF NOT EXISTS orderlines (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  orderid character varying(36),
  variantid character varying(36),
  productname character varying(386),
  variantname character varying(255),
  translatedproductname character varying(386),
  translatedvariantname character varying(255),
  productsku character varying(255),
  productvariantid character varying(255),
  isshippingrequired boolean,
  isgiftcard boolean,
  quantity integer,
  quantityfulfilled integer,
  currency character varying(3),
  unitdiscountamount double precision,
  unitdiscounttype character varying(10),
  unitdiscountreason text,
  unitpricenetamount double precision,
  unitdiscountvalue double precision,
  unitpricegrossamount double precision,
  totalpricenetamount double precision,
  totalpricegrossamount double precision,
  undiscountedunitpricegrossamount double precision,
  undiscountedunitpricenetamount double precision,
  undiscountedtotalpricegrossamount double precision,
  undiscountedtotalpricenetamount double precision,
  taxrate double precision,
  allocations text
);

CREATE INDEX idx_order_lines_product_name_lower_textpattern ON orderlines USING btree (lower((productname)::text) text_pattern_ops);

CREATE INDEX idx_order_lines_translated_product_name ON orderlines USING btree (translatedproductname);

CREATE INDEX idx_order_lines_translated_variant_name ON orderlines USING btree (translatedvariantname);

CREATE INDEX idx_order_lines_variant_name ON orderlines USING btree (variantname);

CREATE INDEX idx_order_lines_variant_name_lower_textpattern ON orderlines USING btree (lower((variantname)::text) text_pattern_ops);

CREATE INDEX idx_order_lines_product_name ON orderlines USING btree (productname);
