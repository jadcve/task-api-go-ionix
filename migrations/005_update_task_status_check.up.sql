DO $$
DECLARE
    rec record;
BEGIN
    FOR rec IN
        SELECT con.conname
        FROM pg_constraint con
        INNER JOIN pg_class rel ON rel.oid = con.conrelid
        INNER JOIN pg_namespace nsp ON nsp.oid = con.connamespace
        WHERE rel.relname = 'tasks'
          AND nsp.nspname = current_schema()
          AND con.contype = 'c'
          AND pg_get_constraintdef(con.oid) ILIKE '%status%'
    LOOP
        EXECUTE format('ALTER TABLE tasks DROP CONSTRAINT IF EXISTS %I', rec.conname);
    END LOOP;
END $$;

ALTER TABLE tasks ALTER COLUMN status SET DEFAULT 'ASSIGNED';

ALTER TABLE tasks
ADD CONSTRAINT chk_tasks_status CHECK (
    status IN (
        'ASSIGNED',
        'STARTED',
        'WAITING',
        'COMPLETED_SUCCESS',
        'COMPLETED_ERROR'
    )
);