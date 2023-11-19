CREATE TABLE IF NOT EXISTS export_events (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  date bigint NOT NULL,
  type character varying(255) NOT NULL,
  parameters text,
  export_file_id uuid NOT NULL,
  user_id uuid
);