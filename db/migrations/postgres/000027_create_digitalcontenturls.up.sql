CREATE TABLE IF NOT EXISTS digital_content_urls (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    token uuid NOT NULL,
    content_id uuid NOT NULL,
    created_at bigint NOT NULL,
    download_num integer NOT NULL,
    line_id uuid
);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_line_id_key UNIQUE (line_id);

ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT digital_content_urls_token_key UNIQUE (token);
