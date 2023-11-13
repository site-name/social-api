CREATE TABLE IF NOT EXISTS users (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  email character varying(128),
  username character varying(64),
  first_name character varying(64),
  last_name character varying(64),
  default_shipping_address_id uuid,
  default_billing_address_id uuid,
  password character varying(128),
  auth_data character varying(128),
  auth_service character varying(32),
  email_verified boolean,
  nickname character varying(64),
  roles character varying(256),
  props jsonb,
  notify_props jsonb,
  last_password_update bigint,
  last_picture_update bigint,
  failed_attempts integer,
  locale character varying(5),
  timezone jsonb,
  mfa_active boolean,
  mfa_secret character varying(128),
  created_at bigint,
  updated_at bigint,
  delete_at bigint,
  is_active boolean,
  note text,
  jwt_token_key text,
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