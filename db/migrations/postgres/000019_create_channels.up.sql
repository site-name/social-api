CREATE TABLE IF NOT EXISTS channels (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(250) NOT NULL,
  is_active boolean NOT NULL,
  slug varchar(255) NOT NULL,
  currency Currency NOT NULL,
  default_country country_code NOT NULL,
  allocation_strategy allocation_strategy NOT NULL DEFAULT 'prioritize_sorting_order',
  order_mark_as_paid_strategy mark_as_paid_strategy NOT NULL DEFAULT 'payment_flow',
  default_transaction_flow_strategy transaction_flow_strategy NOT NULL DEFAULT 'charge',
  automatically_confirm_all_new_orders boolean DEFAULT true,
  allow_unpaid_orders boolean NOT NULL DEFAULT false,
  automatically_fulfill_non_shippable_gift_card boolean DEFAULT true,
  expire_orders_after integer,
  delete_expired_orders_after integer default 60, -- in days
  include_darft_order_is_voucher_usage boolean NOT NULL DEFAULT false,
  automatically_complete_fully_paid_checkouts boolean NOT NULL DEFAULT false,
  annotations jsonb -- This is a JSONB column that will store the annotations of the channel
);

ALTER TABLE ONLY channels
    ADD CONSTRAINT channels_slug_key UNIQUE (slug);

CREATE INDEX IF NOT EXISTS idx_channels_currency ON channels USING btree (currency);

CREATE INDEX IF NOT EXISTS idx_channels_name ON channels USING btree (name);

CREATE INDEX IF NOT EXISTS idx_channels_is_active ON channels USING btree (is_active);

CREATE INDEX IF NOT EXISTS idx_channels_name_lower_textpattern ON channels USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_channels_slug ON channels USING btree (slug);