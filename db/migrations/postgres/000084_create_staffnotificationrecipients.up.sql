CREATE TABLE IF NOT EXISTS staff_notification_recipients (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid,
  staff_email varchar(128),
  active boolean NOT NULL,
  CONSTRAINT staff_email_unique UNIQUE (staff_email)
);