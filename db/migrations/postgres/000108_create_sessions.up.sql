CREATE TABLE IF NOT EXISTS sessions (
  id varchar(36) NOT NULL PRIMARY KEY,
  token VARCHAR(36),
  createat bigint,
  expiresat bigint,
  lastactivityat bigint,
  userid VARCHAR(36),
  deviceid VARCHAR(512),
  roles VARCHAR(64),
  isoauth boolean,
  expirednotify BOOLEAN,
  props VARCHAR(1000)
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (userid);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions (token);
CREATE INDEX IF NOT EXISTS idx_sessions_last_activity_at ON sessions (lastactivityat);
