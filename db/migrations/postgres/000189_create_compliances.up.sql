CREATE TABLE IF NOT EXISTS compliances (
    id varchar(36) NOT NULL PRIMARY KEY,
    created_at bigint NOT NULL,
    user_id varchar(36) NOT NULL,
    status compliance_status NOT NULL,
    count integer NOT NULL,
    "desc" varchar(512) NOT NULL,
    type compliance_type NOT NULL,
    start_at bigint NOT NULL,
    end_at bigint NOT NULL,
    keywords varchar(512) NOT NULL,
    emails varchar(1024) NOT NULL
);

ALTER TABLE ONLY compliances
    ADD CONSTRAINT fk_compliances_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

