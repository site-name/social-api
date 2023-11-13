ALTER TABLE ONLY staff_notification_recipients
    ADD CONSTRAINT fk_staff_notification_recipients_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
