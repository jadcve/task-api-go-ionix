UPDATE tasks
SET
    status = CASE status
        WHEN 'STARTED' THEN 'IN_PROGRESS'
        WHEN 'WAITING' THEN 'IN_PROGRESS'
        WHEN 'COMPLETED_SUCCESS' THEN 'COMPLETED'
        WHEN 'COMPLETED_ERROR' THEN 'CANCELLED'
        ELSE status
    END;

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
        'IN_PROGRESS',
        'COMPLETED',
        'CANCELLED'
    )
);