CREATE TABLE IF NOT EXISTS users (
  id varchar(36) NOT NULL PRIMARY KEY,
  email varchar(128) NOT NULL,
  username varchar(64) NOT NULL,
  first_name varchar(64) NOT NULL,
  last_name varchar(64) NOT NULL,
  default_shipping_address_id varchar(36),
  default_billing_address_id varchar(36),
  password varchar(128) NOT NULL,
  auth_data varchar(128),
  auth_service varchar(32) NOT NULL,
  email_verified boolean NOT NULL,
  nickname varchar(64) NOT NULL,
  roles varchar(256) NOT NULL,
  props jsonb,
  notify_props jsonb,
  last_password_update bigint NOT NULL,
  last_picture_update bigint NOT NULL,
  failed_attempts integer NOT NULL,
  locale varchar(5) NOT NULL,
  timezone jsonb,
  mfa_active boolean NOT NULL,
  mfa_secret varchar(128) NOT NULL,
  created_at bigint NOT NULL,
  updated_at bigint NOT NULL,
  delete_at bigint NOT NULL,
  is_active boolean NOT NULL,
  note text,
  jwt_token_key text NOT NULL,
  last_activity_at bigint NOT NULL,
  terms_of_service_id varchar(36) NOT NULL,
  terms_of_service_created_at bigint NOT NULL,
  disable_welcome_email bool NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_auth_data_key UNIQUE (auth_data);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_key UNIQUE (username);

CREATE INDEX idx_users_all_no_full_name_txt ON users USING gin (to_tsvector('english'::regconfig, (((((username)::text || ' '::text) || (nickname)::text) || ' '::text) || (email)::text)));

CREATE INDEX idx_users_all_txt ON users USING gin (to_tsvector('english'::regconfig, (((((((((username)::text || ' '::text) || (first_name)::text) || ' '::text) || (last_name)::text) || ' '::text) || (nickname)::text) || ' '::text) || (email)::text)));

CREATE INDEX idx_users_email ON users USING btree (email);

CREATE INDEX idx_users_email_lower_text_pattern ON users USING btree (lower((email)::text) text_pattern_ops);

CREATE INDEX idx_users_first_name_lower_text_pattern ON users USING btree (lower((first_name)::text) text_pattern_ops);

CREATE INDEX idx_users_last_name_lower_text_pattern ON users USING btree (lower((last_name)::text) text_pattern_ops);

CREATE INDEX idx_users_metadata ON users USING btree (metadata);

CREATE INDEX idx_users_names_no_full_name_txt ON users USING gin (to_tsvector('english'::regconfig, (((username)::text || ' '::text) || (nickname)::text)));

CREATE INDEX idx_users_names_txt ON users USING gin (to_tsvector('english'::regconfig, (((((((username)::text || ' '::text) || (first_name)::text) || ' '::text) || (last_name)::text) || ' '::text) || (nickname)::text)));

CREATE INDEX idx_users_nickname_lower_text_pattern ON users USING btree (lower((nickname)::text) text_pattern_ops);

CREATE INDEX idx_users_private_metadata ON users USING btree (private_metadata);

CREATE INDEX idx_users_username_lower_text_pattern ON users USING btree (lower((username)::text) text_pattern_ops);