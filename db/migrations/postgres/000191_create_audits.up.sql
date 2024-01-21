CREATE TABLE IF NOT EXISTS audits (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at bigint NOT NULL,
    user_id uuid NOT NULL,
    action VARCHAR(512) NOT NULL,
    extra_info VARCHAR(1024) NOT NULL,
    ip_address VARCHAR(64) NOT NULL,
    session_id uuid NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_audits_user_id ON audits (user_id);
