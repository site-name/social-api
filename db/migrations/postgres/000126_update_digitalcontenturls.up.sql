ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT fk_digital_content_urls_digital_contents FOREIGN KEY (contentid) REFERENCES digital_contents(id) ON DELETE CASCADE;
ALTER TABLE ONLY digital_content_urls
    ADD CONSTRAINT fk_digital_content_urls_order_lines FOREIGN KEY (lineid) REFERENCES order_lines(id) ON DELETE CASCADE;
