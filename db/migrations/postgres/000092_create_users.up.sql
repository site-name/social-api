CREATE TABLE IF NOT EXISTS users (
  id character varying(36) NOT NULL PRIMARY KEY,
  email character varying(128),
  username character varying(64),
  firstname character varying(64),
  lastname character varying(64),
  defaultshippingaddressid character varying(36),
  defaultbillingaddressid character varying(36),
  password character varying(128),
  authdata character varying(128),
  authservice character varying(32),
  emailverified boolean,
  nickname character varying(64),
  roles character varying(256),
  props character varying(4000),
  notifyprops character varying(2000),
  lastpasswordupdate bigint,
  lastpictureupdate bigint,
  failedattempts integer,
  locale character varying(5),
  timezone character varying(256),
  mfaactive boolean,
  mfasecret character varying(128),
  createat bigint,
  updateat bigint,
  deleteat bigint,
  isactive boolean,
  note text,
  jwttokenkey text,
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_authdata_key UNIQUE (authdata);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_key UNIQUE (username);

CREATE INDEX idx_users_all_no_full_name_txt ON users USING gin (to_tsvector('english'::regconfig, (((((username)::text || ' '::text) || (nickname)::text) || ' '::text) || (email)::text)));

CREATE INDEX idx_users_all_txt ON users USING gin (to_tsvector('english'::regconfig, (((((((((username)::text || ' '::text) || (firstname)::text) || ' '::text) || (lastname)::text) || ' '::text) || (nickname)::text) || ' '::text) || (email)::text)));

CREATE INDEX idx_users_email ON users USING btree (email);

CREATE INDEX idx_users_email_lower_textpattern ON users USING btree (lower((email)::text) text_pattern_ops);

CREATE INDEX idx_users_firstname_lower_textpattern ON users USING btree (lower((firstname)::text) text_pattern_ops);

CREATE INDEX idx_users_lastname_lower_textpattern ON users USING btree (lower((lastname)::text) text_pattern_ops);

CREATE INDEX idx_users_metadata ON users USING btree (metadata);

CREATE INDEX idx_users_names_no_full_name_txt ON users USING gin (to_tsvector('english'::regconfig, (((username)::text || ' '::text) || (nickname)::text)));

CREATE INDEX idx_users_names_txt ON users USING gin (to_tsvector('english'::regconfig, (((((((username)::text || ' '::text) || (firstname)::text) || ' '::text) || (lastname)::text) || ' '::text) || (nickname)::text)));

CREATE INDEX idx_users_nickname_lower_textpattern ON users USING btree (lower((nickname)::text) text_pattern_ops);

CREATE INDEX idx_users_private_metadata ON users USING btree (privatemetadata);

CREATE INDEX idx_users_username_lower_textpattern ON users USING btree (lower((username)::text) text_pattern_ops);
