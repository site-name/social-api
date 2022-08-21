CREATE TABLE IF NOT EXISTS customernotes (
  id character varying(36),
	userid character varying(36),
	date bigint,
	content text,
	ispublic boolean,
	customerid character varying(36)
);

CREATE INDEX idx_customer_notes_date ON public.customernotes USING btree (date);

ALTER TABLE ONLY public.customernotes
    ADD CONSTRAINT fk_customernotes_users FOREIGN KEY (userid) REFERENCES public.users(id);