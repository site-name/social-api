CREATE TABLE IF NOT EXISTS staff_notification_recipients (
  id varchar(36) NOT NULL PRIMARY KEY,
  user_id varchar(36),
  staff_email varchar(128),
  active boolean NOT NULL,
  CONSTRAINT staff_email_unique UNIQUE (staff_email)
);