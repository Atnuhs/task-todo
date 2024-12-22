CREATE TYPE task_status AS ENUM('pending', 'doing', 'completed', 'cancelled');

CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    task_id UUID UNIQUE NOT NULL,
    task_name text NOT NULL,
    status task_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
)