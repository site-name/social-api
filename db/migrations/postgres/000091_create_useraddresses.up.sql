CREATE TABLE IF NOT EXISTS user_addresses (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid,
  address_id character varying(36)
);

ALTER TABLE ONLY user_addresses
    ADD CONSTRAINT user_addresses_user_id_address_id_key UNIQUE (user_id, address_id);