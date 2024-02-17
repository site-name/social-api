CREATE TABLE IF NOT EXISTS export_events (
  id varchar(36) NOT NULL PRIMARY KEY,
  date bigint NOT NULL,
  type export_event_type NOT NULL,
  parameters text,
  export_file_id varchar(36) NOT NULL,
  user_id varchar(36)
);