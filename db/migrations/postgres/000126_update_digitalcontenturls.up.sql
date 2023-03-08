ALTER TABLE ONLY digitalcontenturls
    ADD CONSTRAINT fk_digitalcontenturls_digitalcontents FOREIGN KEY (contentid) REFERENCES digitalcontents(id) ON DELETE CASCADE;
ALTER TABLE ONLY digitalcontenturls
    ADD CONSTRAINT fk_digitalcontenturls_orderlines FOREIGN KEY (lineid) REFERENCES orderlines(id) ON DELETE CASCADE;
