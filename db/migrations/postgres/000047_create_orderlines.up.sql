CREATE TABLE IF NOT EXISTS order_lines (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  order_id varchar(36) NOT NULL,
  variant_id varchar(36),
  product_name varchar(386) NOT NULL,
  variant_name varchar(255) NOT NULL,
  translated_product_name varchar(386) NOT NULL,
  translated_variant_name varchar(255) NOT NULL,
  product_sku varchar(255),
  product_variant_id varchar(255),
  is_shipping_required boolean NOT NULL,
  is_giftcard boolean NOT NULL,
  is_gift boolean NOT NULL default false,
  quantity integer NOT NULL,
  quantity_fulfilled integer NOT NULL,
  currency Currency NOT NULL,
  unit_discount_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  unit_discount_type discount_value_type NOT NULL,
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
  base_unit_price_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  undiscounted_base_unit_price_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  tax_rate decimal(5,4),
  tax_class_id varchar(36),
  tax_class_name varchar(255),
  tax_class_private_metadata jsonb,
  tax_class_metadata jsonb,
  is_price_overridden boolean,
  voucher_code varchar(255),
  sale_id varchar(36)
);

CREATE INDEX idx_order_lines_product_name_lower_textpattern ON order_lines USING btree (lower((product_name)::text) text_pattern_ops);

CREATE INDEX idx_order_lines_translated_product_name ON order_lines USING btree (translated_product_name);

CREATE INDEX idx_order_lines_translated_variant_name ON order_lines USING btree (translated_variant_name);

CREATE INDEX idx_order_lines_variant_name ON order_lines USING btree (variant_name);

CREATE INDEX idx_order_lines_variant_name_lower_textpattern ON order_lines USING btree (lower((variant_name)::text) text_pattern_ops);

CREATE INDEX idx_order_lines_product_name ON order_lines USING btree (product_name);