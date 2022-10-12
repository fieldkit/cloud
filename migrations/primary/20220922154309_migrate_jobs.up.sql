CREATE TABLE IF NOT EXISTS fieldkit.gue_jobs
(
    job_id      BIGSERIAL   NOT NULL PRIMARY KEY,
    priority    SMALLINT    NOT NULL,
    run_at      TIMESTAMPTZ NOT NULL,
    job_type    TEXT        NOT NULL,
    args        JSON        NOT NULL,
    error_count INTEGER     NOT NULL DEFAULT 0,
    last_error  TEXT,
    queue       TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_gue_jobs_selector ON fieldkit.gue_jobs (queue, run_at, priority);

COMMENT ON TABLE gue_jobs IS '1';