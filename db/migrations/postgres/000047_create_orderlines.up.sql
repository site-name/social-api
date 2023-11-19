CREATE TABLE IF NOT EXISTS order_lines (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  order_id uuid NOT NULL,
  variant_id uuid,
  product_name character varying(386) NOT NULL,
  variant_name character varying(255) NOT NULL,
  translated_product_name character varying(386) NOT NULL,
  translated_variant_name character varying(255) NOT NULL,
  product_sku character varying(255),
  product_variant_id character varying(255),
  is_shipping_required boolean NOT NULL,
  is_giftcard boolean NOT NULL,
  quantity integer NOT NULL,
  quantity_fulfilled integer NOT NULL,
  currency character varying(3) NOT NULL,
  unit_discount_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  unit_discount_type character varying(10) NOT NULL,
  unit_discount_reason text,
  unit_price_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  unit_discount_value decimal(12,3) NOT NULL DEFAULT 0.00,
  unit_price_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_price_net_amount decimal(12,3),
  total_price_gross_amount decimal(12,3),
  undiscounted_unit_price_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_unit_price_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_total_price_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_total_price_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  tax_rate decimal(5,4)
);

CREATE INDEX idx_order_lines_product_name_lower_textpattern ON order_lines USING btree (lower((product_name)::text) text_pattern_ops);

CREATE INDEX idx_order_lines_translated_product_name ON order_lines USING btree (translated_product_name);

CREATE INDEX idx_order_lines_translated_variant_name ON order_lines USING btree (translated_variant_name);

CREATE INDEX idx_order_lines_variant_name ON order_lines USING btree (variant_name);

CREATE INDEX idx_order_lines_variant_name_lower_textpattern ON order_lines USING btree (lower((variant_name)::text) text_pattern_ops);

CREATE INDEX idx_order_lines_product_name ON order_lines USING btree (product_name);