CREATE TABLE IF NOT EXISTS exportevents (
  id character varying(36) NOT NULL PRIMARY KEY,
  date bigint,
  type character varying(255),
  parameters text,
  exportfileid character varying(36),
  userid character varying(36)
);

ALTER TABLE ONLY exportevents
    ADD CONSTRAINT fk_exportevents_exportfiles FOREIGN KEY (exportfileid) REFERENCES exportfiles(id);

ALTER TABLE ONLY exportevents
    ADD CONSTRAINT fk_exportevents_users FOREIGN KEY (userid) REFERENCES users(id);

