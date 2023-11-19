CREATE TABLE IF NOT EXISTS jobs (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  type character varying(32) NOT NULL,
  priority bigint NOT NULL,
  created_at bigint NOT NULL,
  start_at bigint NOT NULL,
  last_activity_at bigint NOT NULL,
  status character varying(32) NOT NULL,
  progress bigint NOT NULL,
  data jsonb
);

CREATE INDEX idx_jobs_type ON jobs USING btree (type);