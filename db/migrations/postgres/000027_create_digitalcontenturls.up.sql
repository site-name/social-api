CREATE TABLE IF NOT EXISTS digital_content_urls (
    id varchar(36) NOT NULL PRIMARY KEY,
    token varchar(36) NOT NULL,
    content_id varchar(36) NOT NULL,
    created_at bigint NOT NULL,
    download_num integer NOT NULL,
    line_id varchar(36)
);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_line_id_key UNIQUE (line_id);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_token_key UNIQUE (token);
