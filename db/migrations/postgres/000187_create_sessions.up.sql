CREATE TABLE IF NOT EXISTS sessions (
  id VARCHAR(36) PRIMARY KEY,
  token VARCHAR(36),
  created_at bigint,
  expires_at bigint,
  last_activity_at bigint,
  user_id VARCHAR(36),
  device_id VARCHAR(512),
  roles VARCHAR(64),
  is_oauth BOOLEAN,
  expired_notify BOOLEAN,
  props jsonb,
  local BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions (token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_create_at ON sessions (create_at);
CREATE INDEX IF NOT EXISTS idx_sessions_last_activity_at ON sessions (last_activity_at);