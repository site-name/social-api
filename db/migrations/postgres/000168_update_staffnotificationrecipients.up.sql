ALTER TABLE ONLY staff_notification_recipients
    ADD CONSTRAINT fk_staff_notification_recipients_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;
