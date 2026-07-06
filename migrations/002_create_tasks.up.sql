CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    due_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ASSIGNED' CHECK (
        status IN (
            'ASSIGNED',
            'IN_PROGRESS',
            'COMPLETED',
            'CANCELLED'
        )
    ),
    assigned_to BIGINT NOT NULL REFERENCES users (id),
    created_by BIGINT NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks (assigned_to);

CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON tasks (created_by);