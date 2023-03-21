CREATE TABLE IF NOT EXISTS sessions (
  id VARCHAR(36) PRIMARY KEY,
  token VARCHAR(36),
  createat bigint,
  expiresat bigint,
  lastactivityat bigint,
  userid VARCHAR(36),
  devideid VARCHAR(512),
  roles VARCHAR(64),
  isoauth BOOLEAN,
  expirednotify BOOLEAN,
  props jsonb,
  local BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (userid);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions (token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expiresat);
CREATE INDEX IF NOT EXISTS idx_sessions_create_at ON sessions (createat);
CREATE INDEX IF NOT EXISTS idx_sessions_last_activity_at ON sessions (lastactivityat);
