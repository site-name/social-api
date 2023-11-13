CREATE TABLE IF NOT EXISTS staff_notification_recipients (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid,
  staff_email character varying(128),
  active boolean
);