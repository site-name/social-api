CREATE TABLE IF NOT EXISTS export_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  date bigint,
  type character varying(255),
  parameters text,
  export_file_id uuid,
  user_id character varying(36)
);