CREATE TABLE IF NOT EXISTS shop_staffs (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  staff_id uuid NOT NULL,
  created_at bigint NOT NULL,
  end_at bigint,
  salary_period varchar(10) NOT NULL,
  salary decimal(12,3) NOT NULL DEFAULT 0.00,
  salary_currency varchar(3) NOT NULL
);
