CREATE TABLE IF NOT EXISTS shops (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  name varchar(100) NOT NULL,
  description varchar(200) NOT NULL,
  top_menu_id varchar(36),
  include_taxes_in_price boolean,
  display_gross_prices boolean,
  charge_taxes_on_shipping boolean,
  track_inventory_by_default boolean,
  default_weight_unit varchar(10) NOT NULL,
  automatic_fulfillment_digital_products boolean,
  default_digital_max_downloads integer,
  default_digital_url_valid_days integer,
  address_id varchar(36),
  default_mail_sender_name varchar(78) NOT NULL,
  default_mail_sender_address text NOT NULL,
  customer_set_password_url text,
  automatically_confirm_all_new_orders boolean,
  fulfillment_auto_approve boolean,
  fulfillment_allow_unpaid boolean,
  gift_card_expiry_type varchar(32) NOT NULL,
  gift_card_expiry_period_type varchar(32) NOT NULL,
  gift_card_expiry_period integer,
  automatically_fulfill_non_shippable_gift_card boolean
);

CREATE INDEX idx_shops_description ON shops USING btree (description);

CREATE INDEX idx_shops_description_lower_text_pattern ON shops USING btree (lower((description)::text) text_pattern_ops);

CREATE INDEX idx_shops_name ON shops USING btree (name);

CREATE INDEX idx_shops_name_lower_text_pattern ON shops USING btree (lower((name)::text) text_pattern_ops);