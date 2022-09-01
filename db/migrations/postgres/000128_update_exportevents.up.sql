ALTER TABLE ONLY exportevents
    ADD CONSTRAINT fk_exportevents_exportfiles FOREIGN KEY (exportfileid) REFERENCES exportfiles(id);
ALTER TABLE ONLY exportevents
    ADD CONSTRAINT fk_exportevents_users FOREIGN KEY (userid) REFERENCES users(id);
