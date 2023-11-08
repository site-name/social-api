CREATE TABLE IF NOT EXISTS export_events (
  id character varying(36) NOT NULL PRIMARY KEY,
  date bigint,
  type character varying(255),
  parameters text,
  exportfileid character varying(36),
  userid character varying(36)
);

