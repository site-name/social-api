CREATE TABLE IF NOT EXISTS order_lines (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  order_id uuid,
  variant_id uuid,
  product_name character varying(386),
  variant_name character varying(255),
  translated_product_name character varying(386),
  translated_variant_name character varying(255),
  product_sku character varying(255),
  product_variant_id character varying(255),
  is_shipping_required boolean,
  is_giftcard boolean,
  quantity integer,
  quantity_fulfilled integer,
  currency character varying(3),
  unit_discount_amount double precision,
  unit_discount_type character varying(10),
  unit_discount_reason text,
  unit_price_net_amount double precision,
  unit_discount_value double precision,
  unit_price_gross_amount double precision,
  total_price_net_amount double precision,
  total_price_gross_amount double precision,
  undiscounted_unit_price_gross_amount double precision,
  undiscounted_unit_price_net_amount double precision,
  undiscounted_total_price_gross_amount double precision,
  undiscounted_total_price_net_amount double precision,
  tax_rate double precision
);

CREATE INDEX idx_order_lines_product_name_lower_textpattern ON order_lines USING btree (lower((product_name)::text) text_pattern_ops);

CREATE INDEX idx_order_lines_translated_product_name ON order_lines USING btree (translated_product_name);

CREATE INDEX idx_order_lines_translated_variant_name ON order_lines USING btree (translated_variant_name);

CREATE INDEX idx_order_lines_variant_name ON order_lines USING btree (variant_name);

CREATE INDEX idx_order_lines_variant_name_lower_textpattern ON order_lines USING btree (lower((variant_name)::text) text_pattern_ops);

CREATE INDEX idx_order_lines_product_name ON order_lines USING btree (product_name);