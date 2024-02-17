CREATE TABLE IF NOT EXISTS jobs (
  id varchar(36) NOT NULL PRIMARY KEY,
  type job_type NOT NULL,
  priority bigint NOT NULL,
  created_at bigint NOT NULL,
  start_at bigint NOT NULL,
  last_activity_at bigint NOT NULL,
  status job_status NOT NULL,
  progress bigint NOT NULL,
  data jsonb
);

CREATE INDEX idx_jobs_type ON jobs USING btree (type);