CREATE TABLE IF NOT EXISTS jobs (
  id character varying(36) NOT NULL PRIMARY KEY,
  type character varying(32),
  priority bigint,
  createat bigint,
  startat bigint,
  lastactivityat bigint,
  status character varying(32),
  progress bigint,
  data jsonb
);

CREATE INDEX idx_jobs_type ON jobs USING btree (type);
