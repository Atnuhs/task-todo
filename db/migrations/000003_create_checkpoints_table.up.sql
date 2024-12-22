CREATE TABLE IF NOT EXISTS checkpoints (
    projection_id VARCHAR(255) PRIMARY KEY,
    last_checkpoint INT NOT NULL
);

INSERT INTO checkpoints (projection_id, last_checkpoint)
VALUES ('task_projection', 0)
ON CONFLICT (projection_id) DO NOTHING;
