CREATE TABLE IF NOT EXISTS shop_staffs (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  staff_id uuid,
  created_at bigint,
  end_at bigint,
  salary_period character varying(10),
  salary double precision,
  salary_currency character varying(5)
);
