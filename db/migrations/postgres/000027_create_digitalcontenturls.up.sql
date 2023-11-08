CREATE TABLE IF NOT EXISTS digital_content_urls (
    id character varying(36) NOT NULL,
    token character varying(36),
    contentid character varying(36),
    createat bigint,
    downloadnum integer,
    lineid character varying(36)
);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_lineid_key UNIQUE (lineid);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_token_key UNIQUE (token);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_pkey PRIMARY KEY (id);
