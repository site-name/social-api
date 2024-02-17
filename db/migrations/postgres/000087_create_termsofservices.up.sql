CREATE TABLE IF NOT EXISTS terms_of_services (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  user_id varchar(36) NOT NULL,
  text varchar(65535) NOT NULL
);