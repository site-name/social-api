CREATE TABLE IF NOT EXISTS sessions (
  id varchar(36) NOT NULL PRIMARY KEY,
  token varchar(36) NOT NULL,
  created_at bigint NOT NULL,
  expires_at bigint NOT NULL,
  last_activity_at bigint NOT NULL,
  user_id varchar(36) NOT NULL,
  device_id VARCHAR(512) NOT NULL,
  roles VARCHAR(64) NOT NULL,
  is_oauth BOOLEAN NOT NULL,
  expired_notify BOOLEAN NOT NULL,
  props jsonb,
  local boolean NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions (token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_create_at ON sessions (created_at);
CREATE INDEX IF NOT EXISTS idx_sessions_last_activity_at ON sessions (last_activity_at);