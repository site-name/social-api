CREATE TABLE IF NOT EXISTS vouchercustomers (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  customeremail character varying(128)
);

ALTER TABLE ONLY public.vouchercustomers
    ADD CONSTRAINT vouchercustomers_voucherid_customeremail_key UNIQUE (voucherid, customeremail);

ALTER TABLE ONLY public.vouchercustomers
    ADD CONSTRAINT fk_vouchercustomers_vouchers FOREIGN KEY (voucherid) REFERENCES public.vouchers(id) ON DELETE CASCADE;
