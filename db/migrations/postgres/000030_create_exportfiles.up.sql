CREATE TABLE IF NOT EXISTS export_files (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid,
  content_file text,
  created_at bigint,
  updated_at bigint
);