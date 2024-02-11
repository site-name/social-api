
DO $$
BEGIN
    IF NOT EXISTS (
       SELECT conname
       FROM pg_constraint
       WHERE conname = 'addresses_user_id_fk'
    )
    THEN
        ALTER TABLE addresses ADD CONSTRAINT addresses_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;
