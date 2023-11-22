CREATE TABLE IF NOT EXISTS compliances (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at bigint NOT NULL,
    user_id uuid NOT NULL,
    status ComplianceStatus NOT NULL,
    count integer NOT NULL,
    "desc" varchar(512) NOT NULL,
    type ComplianceType NOT NULL,
    start_at bigint NOT NULL,
    end_at bigint NOT NULL,
    keywords varchar(512) NOT NULL,
    emails varchar(1024) NOT NULL
);

ALTER TABLE ONLY compliances
    ADD CONSTRAINT fk_compliances_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

