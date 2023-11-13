CREATE TABLE IF NOT EXISTS jobs (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  type character varying(32),
  priority bigint,
  created_at bigint,
  start_at bigint,
  last_activity_at bigint,
  status character varying(32),
  progress bigint,
  data jsonb
);

CREATE INDEX idx_jobs_type ON jobs USING btree (type);