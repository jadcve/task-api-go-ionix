CREATE TABLE IF NOT EXISTS task_comments (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users (id),
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_task_comments_task_id ON task_comments (task_id);

CREATE INDEX IF NOT EXISTS idx_task_comments_user_id ON task_comments (user_id);