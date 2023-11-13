CREATE TABLE IF NOT EXISTS digital_content_urls (
    id character varying(36) NOT NULL,
    token uuid,
    content_id uuid,
    created_at bigint,
    download_num integer,
    line_id character varying(36)
);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_line_id_key UNIQUE (line_id);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_token_key UNIQUE (token);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_pkey PRIMARY KEY (id);