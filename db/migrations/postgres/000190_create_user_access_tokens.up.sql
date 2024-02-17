CREATE TABLE IF NOT EXISTS user_access_tokens (
    id varchar(36) NOT NULL PRIMARY KEY,
    token varchar(36) NOT NULL,
    user_id varchar(36) NOT NULL,
    description varchar(255) NOT NULL,
    is_active boolean NOT NULL DEFAULT true
);

ALTER TABLE ONLY user_access_tokens
    ADD CONSTRAINT fk_user_access_tokens_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
