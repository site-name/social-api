CREATE TABLE IF NOT EXISTS digitalcontenturls (
    id character varying(36) NOT NULL,
    token character varying(36),
    contentid character varying(36),
    createat bigint,
    downloadnum integer,
    lineid character varying(36)
);

ALTER TABLE ONLY digitalcontenturls
    ADD CONSTRAINT digitalcontenturls_lineid_key UNIQUE (lineid);

ALTER TABLE ONLY digitalcontenturls
    ADD CONSTRAINT digitalcontenturls_token_key UNIQUE (token);

ALTER TABLE ONLY digitalcontenturls
    ADD CONSTRAINT fk_digitalcontenturls_digitalcontents FOREIGN KEY (contentid) REFERENCES digitalcontents(id) ON DELETE CASCADE;
    
ALTER TABLE ONLY digitalcontenturls
    ADD CONSTRAINT fk_digitalcontenturls_orderlines FOREIGN KEY (lineid) REFERENCES orderlines(id) ON DELETE CASCADE;

ALTER TABLE ONLY digitalcontenturls
    ADD CONSTRAINT digitalcontenturls_pkey PRIMARY KEY (id);
