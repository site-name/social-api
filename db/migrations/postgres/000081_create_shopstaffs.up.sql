CREATE TABLE IF NOT EXISTS shop_staffs (
  id varchar(36) NOT NULL PRIMARY KEY,
  staff_id varchar(36) NOT NULL,
  created_at bigint NOT NULL,
  end_at bigint,
  salary_period staff_salary_period NOT NULL,
  salary decimal(12,3) NOT NULL DEFAULT 0.00,
  salary_currency Currency NOT NULL
);

CREATE UNIQUE INDEX shop_staff_staff_id_unique_idx ON shop_staffs (staff_id);
