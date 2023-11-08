CREATE TABLE IF NOT EXISTS staff_notification_recipients (
  id character varying(36) NOT NULL PRIMARY KEY,
  userid character varying(36),
  staffemail character varying(128),
  active boolean
);

